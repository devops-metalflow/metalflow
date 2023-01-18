package v1

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	n_websocket "golang.org/x/net/websocket"
	"io"
	"metalflow/pkg/global"
	"metalflow/pkg/request"
	"metalflow/pkg/response"
	"metalflow/pkg/utils"
	"metalflow/pkg/vncproxy"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

var upgrade = websocket.Upgrader{
	// allow cross-domain.
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// PtyRequestMsg Pseudo-terminal pty basic configuration information struct.
type PtyRequestMsg struct {
	Term     string
	Columns  uint32
	Rows     uint32
	Width    uint32
	Height   uint32
	ModeList string
}

type SshClient struct {
	client         *ssh.Client
	channel        ssh.Channel // separate the ssh channel because subsequent pty changes require resizing.
	channelRequest <-chan *ssh.Request
}

type Ssh struct {
	sshClient  *SshClient
	sftpClient *sftp.Client
}

type ClientsInfo struct {
	lock sync.RWMutex
	data map[string]*Ssh
}

var clients = ClientsInfo{
	lock: sync.RWMutex{},
	data: make(map[string]*Ssh),
}

// NodeConnect tests the ssh connection, generate the corresponding ssh and sftp instances,
// and return the unique connection specified id.
func NodeConnect(c *gin.Context) {
	var req request.NodeShellConnectRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("params binding failed, please check the data type")
		return
	}

	client, err := utils.GetSshClient(utils.NewSshConfig(req.Address, int(req.SshPort), req.Username, req.Password))
	if err != nil {
		global.Log.Error(fmt.Sprintf("failed to establish ssh connection：%v", err))
		response.FailWithMsg("unable to establish ssh connection")
		return
	}
	// open the ssh channel channel.
	channel, incomingRequests, err := client.Conn.OpenChannel("session", nil)
	if err != nil {
		global.Log.Error(fmt.Sprintf("failed to establish ssh channel：%v", err))
		return
	}
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		global.Log.Error(fmt.Sprintf("failed to establish sftp connection：%v", err))
		response.FailWithMsg("unable to establish sftp connection")
		return
	}

	sshId := utils.RandString(15) // nolint:gomnd
	newSsh := &Ssh{
		sshClient: &SshClient{
			client:         client,
			channel:        channel,
			channelRequest: incomingRequests,
		},
		sftpClient: sftpClient,
	}
	clients.lock.Lock()
	clients.data[sshId] = newSsh
	clients.lock.Unlock()
	response.SuccessWithData(sshId)
}

// NodeShellWs starts the machine shell connection.
// nolint:funlen
// nolint:gocyclo
func NodeShellWs(c *gin.Context) { //nolint:gocyclo
	var req request.NodeShellWsRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("parameter binding failed, please check the data type")
		return
	}

	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		clients.lock.Lock()
		delete(clients.data, req.SshId)
		clients.lock.Unlock()
		global.Log.Error("upgrade websocket connection failed", err)
		return
	}

	defer func(conn *websocket.Conn) {
		clients.lock.Lock()
		delete(clients.data, req.SshId)
		clients.lock.Unlock()
		err = conn.Close()
		if err != nil {
			global.Log.Error("failed to close websocket connection", err)
			return
		}
	}(conn)

	// establish connection.
	clients.lock.RLock()
	cli, ok := clients.data[req.SshId]
	clients.lock.RUnlock()
	if !ok {
		global.Log.Error(fmt.Sprintf("failed to establish ssh connection：%v", err))
		_ = conn.WriteMessage(websocket.TextMessage, []byte("\n"+err.Error()))
		return
	}

	defer func(client *ssh.Client) {
		err = client.Close()
		if err != nil {
			global.Log.Error(fmt.Sprintf("failed to close ssh client：%v", err))
		}
	}(cli.sshClient.client)
	defer func(channel ssh.Channel) {
		err = channel.Close()
		if err != nil {
			global.Log.Error(fmt.Sprintf("failed to close ssh channel：%v", err))
		}
	}(cli.sshClient.channel)

	// handle requests that require a reply.
	go func() {
		for r := range cli.sshClient.channelRequest {
			if r.WantReply {
				_ = r.Reply(false, nil)
			}
		}
	}()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	var modeList []byte
	for k, v := range modes {
		kv := struct {
			Key byte
			Val uint32
		}{k, v}
		modeList = append(modeList, ssh.Marshal(&kv)...)
	}
	modeList = append(modeList, 0)

	// send pty
	var rows = uint32(25) //nolint:gomnd
	var cols = uint32(80) //nolint:gomnd
	if req.Rows != 0 {
		rows = uint32(req.Rows)
	}
	if req.Cols != 0 {
		cols = uint32(req.Cols)
	}

	ptyReq := PtyRequestMsg{
		Term:     "xterm",
		Columns:  cols,
		Rows:     rows,
		Width:    rows,
		Height:   cols,
		ModeList: string(modeList),
	}
	ok, err = cli.sshClient.channel.SendRequest("pty-req", true, ssh.Marshal(&ptyReq))
	if !ok || err != nil {
		global.Log.Error(fmt.Sprintf("send pty failed：%v", err))
		_ = conn.WriteMessage(websocket.TextMessage, []byte("\n"+err.Error()))
		return
	}

	// send shell.
	ok, err = cli.sshClient.channel.SendRequest("shell", true, nil)
	if !ok || err != nil {
		global.Log.Error(fmt.Sprintf("Send shell failed: %v", err))
		_ = conn.WriteMessage(websocket.TextMessage, []byte("\n"+err.Error()))
		return
	}

	// handle data reading and writing. That is, read the returned command from the remote host to buf,
	// and then write it back to the websocket conn from buf.
	go func() {
		br := bufio.NewReader(cli.sshClient.channel)
		var buf []byte

		// read from buf and write back to websocket every 100 microseconds.
		t := time.NewTimer(time.Millisecond * 100) // nolint:gomnd
		defer t.Stop()
		r := make(chan rune)

		go func() {
			for {
				// continuously get the output from the channel and put it into r.
				x, size, err := br.ReadRune() //nolint:govet
				if err != nil {
					global.Log.Warn(fmt.Sprintf("read shell warnings: %v", err))
					break
				}
				if size > 0 {
					r <- x
				}
			}
		}()

		for {
			select {
			case <-t.C:
				if len(buf) != 0 {
					err = conn.WriteMessage(websocket.TextMessage, buf)
					buf = []byte{}
					if err != nil {
						global.Log.Error(fmt.Sprintf("write data to %s failed: %v", conn.RemoteAddr(), err))
						return
					}
				}
				t.Reset(time.Millisecond * 100) // nolint:gomnd
			case d := <-r:
				if d != utf8.RuneError {
					p := make([]byte, utf8.RuneLen(d))
					utf8.EncodeRune(p, d)
					buf = append(buf, p...)
				} else {
					buf = append(buf, []byte("@")...)
				}
			}
		}
	}()

	active := time.Now()
	// timeout processing.
	go func() {
		for {
			// check every 5 minutes if the user has not entered data.
			timer := time.NewTimer(5 * time.Second) // nolint:gomnd
			<-timer.C

			cost := time.Since(active)
			if cost.Minutes() >= 120 { //nolint:gomnd
				// if there is no activity for more than 120 minutes, the connection will be automatically closed.
				_ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("\r\ninactive for more than [%s],"+
					" automatically disconnected", cost.String())))
				_ = conn.Close()
				_ = timer.Stop()
				break
			}
		}
	}()
	_ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("\r\nterminal %s successfully connected", req.Address)))
	_ = conn.WriteMessage(websocket.TextMessage, []byte("\r\ntip: if there is no activity for more than two hours,"+
		" it will be automatically disconnected\r\n\r\n"))

	// Continue to read the commands entered by the user from the websocket connection and pass them to the channel of the remote host.
	for {
		active = time.Now()
		// read the data in websocket, p is the command entered by the user.
		message, p, err := conn.ReadMessage()
		if err != nil {
			global.Log.Warn(fmt.Sprintf("connection %s has been disconnected", conn.RemoteAddr()))
			break
		}

		if message == websocket.TextMessage {
			cmd := string(p)
			err = utils.IsSafetyCmd(cmd)
			if err != nil {
				_ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("\r\n\r\n%s\r\n\r\n", err.Error())))
				continue
			}
			_, err = cli.sshClient.channel.Write(p)
			if err != nil {
				break
			}
		}
	}
}

// RFC 4254 Section 6.7.
type ptyWindowChangeMsg struct {
	Columns uint32
	Rows    uint32
	Width   uint32
	Height  uint32
}

// ResizeWs adjusts the size of the terminal.
func ResizeWs(c *gin.Context) {
	var req request.ResizeWsStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("params binding failed, please check the data type")
		return
	}

	clients.lock.RLock()
	cli, ok := clients.data[req.SshId]
	clients.lock.RUnlock()
	if !ok {
		response.FailWithMsg("unable to find ssh instance to connect")
		return
	}

	sshReq := ptyWindowChangeMsg{
		Columns: uint32(req.Width),
		Rows:    uint32(req.High),
		Width:   uint32(req.Width * 8), //nolint:gomnd
		Height:  uint32(req.High * 8),  //nolint:gomnd
	}
	_, err = cli.sshClient.channel.SendRequest("window-change", false, ssh.Marshal(&sshReq))
	if err != nil {
		response.FailWithMsg("failed to resize terminal")
	}
}

// GetSshDirInfo get all file information under the folder path.
func GetSshDirInfo(c *gin.Context) {
	var req request.NodeShellFileStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("params binding failed, please check the data type")
		return
	}

	clients.lock.RLock()
	cli, ok := clients.data[req.SshId]
	clients.lock.RUnlock()
	if !ok {
		response.FailWithMsg("unable to find ssh instance to connect to")
		return
	}
	if err != nil {
		global.Log.Error("sftp client create failed： ", err)
		return
	}
	fileInfos, err := cli.sftpClient.ReadDir(req.Path)
	if err != nil {
		global.Log.Error("sftp get files failed: ", err)
		return
	}
	var fileCount, dirCount uint
	var fileList []map[string]string
	for _, info := range fileInfos {
		fileInfo := map[string]string{}
		fileInfo["path"] = path.Join(req.Path, info.Name())
		fileInfo["name"] = info.Name()
		fileInfo["mode"] = info.Mode().String()
		fileInfo["size"] = utils.ByteSize(info.Size())
		fileInfo["mod_time"] = info.ModTime().Format("2006-01-02 15:04:05")
		if info.IsDir() {
			fileInfo["type"] = "d"
			dirCount += 1
		} else {
			fileInfo["type"] = "f"
			fileCount += 1
		}
		fileList = append(fileList, fileInfo)
	}
	sort.Slice(fileList, func(i, j int) bool {
		return fileList[i]["type"] < fileList[j]["type"]
	})

	// internal method, processing path information.
	pathHandler := func(dirPath string) (paths []map[string]string) {
		tmp := strings.Split(dirPath, "/")

		var dirs []string
		if strings.HasPrefix(dirPath, "/") {
			dirs = append(dirs, "/")
		}

		for _, item := range tmp {
			name := strings.TrimSpace(item)
			if len(name) > 0 {
				dirs = append(dirs, name)
			}
		}

		for i, item := range dirs {
			fullPath := path.Join(dirs[:i+1]...)
			pathInfo := map[string]string{}
			pathInfo["name"] = item
			pathInfo["dir"] = fullPath
			paths = append(paths, pathInfo)
		}
		return paths
	}

	resp := response.ShellWsFilesResponseStruct{
		Files:      fileList,
		FileCount:  fileCount,
		Paths:      pathHandler(req.Path),
		DirCount:   dirCount,
		CurrentDir: req.Path,
	}
	response.SuccessWithData(resp)
}

// GetSshFile read file content.
func GetSshFile(c *gin.Context) {
	var req request.NodeShellFileStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("Parameter binding failed, please check the data type")
		return
	}
	clients.lock.RLock()
	cli, ok := clients.data[req.SshId]
	clients.lock.RUnlock()
	if !ok {
		response.FailWithMsg("")
		return
	}
	file, err := cli.sftpClient.Open(req.Path)
	if err != nil {
		global.Log.Error("sftp read file failed: ", err)
		return
	}
	defer func(file *sftp.File) {
		_ = file.Close()
	}(file)
	all, err := io.ReadAll(file)
	if err != nil {
		return
	}
	response.SuccessWithData(string(all))
}

// DownloadFile download files via sftp.
func DownloadFile(c *gin.Context) {
	var req request.NodeShellFileStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("Parameter binding failed, please check the data type")
		return
	}

	clients.lock.RLock()
	cli, ok := clients.data[req.SshId]
	clients.lock.RUnlock()
	if !ok {
		response.FailWithMsg("unable to find ssh instance to connect")
		return
	}
	file, err := cli.sftpClient.Open(req.Path)
	if err != nil {
		return
	}
	defer func(file *sftp.File) {
		_ = file.Close()
	}(file)
	_, _ = io.Copy(c.Writer, file)
}

// UpdateFile update remote file content through sftp protocol.
func UpdateFile(c *gin.Context) {
	var req request.ModifyFileRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("Parameter binding failed, please check the data type")
		return
	}

	clients.lock.RLock()
	cli, ok := clients.data[req.SshId]
	clients.lock.RUnlock()
	if !ok {
		response.FailWithMsg("unable to find ssh instance to connect")
		return
	}
	// determine whether the file has read and write permissions.
	file, err := cli.sftpClient.OpenFile(req.Path, os.O_RDWR|os.O_TRUNC)
	if err != nil {
		response.FailWithMsg(fmt.Sprintf("file %s does not have read and write permissions", req.Path))
		return
	}
	defer file.Close()
	_, err = file.WriteAt([]byte(req.Content), 0)
	if err != nil {
		response.FailWithMsg(fmt.Sprintf("failed to write the content of file %s", file.Name()))
		return
	}
	response.SuccessWithData(fmt.Sprintf("file %s updated successfully", file.Name()))
}

// NodeVncWs execute remote VNC connection. this will upgrade http to websocket.
func NodeVncWs(c *gin.Context) {
	var req request.NodeVncWsRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("parameter binding failed, please check the data type")
		return
	}
	p := vncproxy.New(fmt.Sprintf("%s:%d", req.Address, req.Port))
	handler := n_websocket.Handler(p.ServeWS)
	handler.ServeHTTP(c.Writer, c.Request)
}
