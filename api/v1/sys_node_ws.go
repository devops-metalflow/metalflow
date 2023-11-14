package v1

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	nwebsocket "golang.org/x/net/websocket"
	"io"
	"metalflow/pkg/global"
	"metalflow/pkg/request"
	"metalflow/pkg/response"
	"metalflow/pkg/utils"
	"metalflow/pkg/vncproxy"
	"net/http"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

var upgrade = websocket.Upgrader{
	// 允许跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// PtyRequestMsg 伪终端pty基本配置信息
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
	channel        ssh.Channel // 将channel独立出来，是因为后续pty变化需要重新设置尺寸
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

// NodeConnect 测试ssh连接，生成对应的ssh与sftp实例，返回唯一连接指定id
func NodeConnect(c *gin.Context) {
	var req request.NodeShellConnectRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	client, err := utils.GetSshClient(utils.NewSshConfig(req.Address, int(req.SshPort), req.Username, req.Password))
	if err != nil {
		global.Log.Error(fmt.Sprintf("建立ssh连接失败：%v", err))
		response.FailWithMsg("无法建立ssh连接")
		return
	}
	// 开启ssh通道channel
	channel, incomingRequests, err := client.Conn.OpenChannel("session", nil)
	if err != nil {
		global.Log.Error(fmt.Sprintf("建立ssh通道失败：%v", err))
		return
	}
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		global.Log.Error(fmt.Sprintf("建立sftp连接失败：%v", err))
		response.FailWithMsg("无法建立sftp连接")
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

// NodeShellWs 启动机器shell连接
//
//nolint:funlen
//nolint:gocyclo
func NodeShellWs(c *gin.Context) { //nolint:gocyclo
	var req request.NodeShellWsRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		clients.lock.Lock()
		delete(clients.data, req.SshId)
		clients.lock.Unlock()
		global.Log.Error("升级websocket连接失败", err)
		return
	}

	defer func(conn *websocket.Conn) {
		clients.lock.Lock()
		delete(clients.data, req.SshId)
		clients.lock.Unlock()
		err = conn.Close()
		if err != nil {
			global.Log.Error("关闭websocket连接失败", err)
			return
		}
	}(conn)

	// 建立连接
	clients.lock.RLock()
	cli, ok := clients.data[req.SshId]
	clients.lock.RUnlock()
	if !ok {
		global.Log.Error(fmt.Sprintf("建立ssh连接失败：%v", err))
		_ = conn.WriteMessage(websocket.TextMessage, []byte("\n"+err.Error()))
		return
	}

	defer func(client *ssh.Client) {
		err = client.Close()
		if err != nil {
			global.Log.Error(fmt.Sprintf("关闭ssh 客户端失败：%v", err))
		}
	}(cli.sshClient.client)
	defer func(channel ssh.Channel) {
		err = channel.Close()
		if err != nil {
			global.Log.Error(fmt.Sprintf("关闭ssh通道失败：%v", err))
		}
	}(cli.sshClient.channel)

	// 处理需要回复的请求
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

	// 发送pty
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
		global.Log.Error(fmt.Sprintf("发送pty失败：%v", err))
		_ = conn.WriteMessage(websocket.TextMessage, []byte("\n"+err.Error()))
		return
	}

	// 发送shell
	ok, err = cli.sshClient.channel.SendRequest("shell", true, nil)
	if !ok || err != nil {
		global.Log.Error(fmt.Sprintf("发送shell失败%v", err))
		_ = conn.WriteMessage(websocket.TextMessage, []byte("\n"+err.Error()))
		return
	}

	// 处理数据读写。即从远程主机中读取返回的命令到buf中，再由buf写回到websocket conn
	go func() {
		br := bufio.NewReader(cli.sshClient.channel)
		var buf []byte

		// 每隔100微妙从buf中读取并写回到websocket
		t := time.NewTimer(time.Millisecond * 100) //nolint:gomnd
		defer t.Stop()
		r := make(chan rune)

		go func() {
			for {
				// 不断从channel中获取输出，并放入到r中
				x, size, err := br.ReadRune() //nolint:govet
				if err != nil {
					global.Log.Warn(fmt.Sprintf("读取shell警告%v", err))
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
					// Process the results and colorize special characters
					retString := string(buf)
					retString = formatLine(retString)
					ret := []byte(retString)
					err = conn.WriteMessage(websocket.TextMessage, ret)
					buf = []byte{}
					if err != nil {
						global.Log.Error(fmt.Sprintf("数据写出到%s失败%v", conn.RemoteAddr(), err))
						return
					}
				}
				t.Reset(time.Millisecond * 100) //nolint:gomnd
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
	// 超时处理
	go func() {
		for {
			// 每5分钟检查一次是否用户没有输入数据
			timer := time.NewTimer(5 * time.Second) //nolint:gomnd
			<-timer.C

			cost := time.Since(active)
			if cost.Minutes() >= 120 { //nolint:gomnd
				// 超时30分钟未活动， 自动关闭连接
				_ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("\r\n已超过【%s】未活动，自动断开连接", cost.String())))
				_ = conn.Close()
				_ = timer.Stop()
				break
			}
		}
	}()

	// 持续从websocket连接中读取用户输入的命令，并将其传递给远程主机的channel中
	for {
		active = time.Now()
		// 读取websocket中数据，p即为用户输入的命令
		message, p, err := conn.ReadMessage()
		if err != nil {
			global.Log.Warn(fmt.Sprintf("连接%s已断开", conn.RemoteAddr()))
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

// ResizeWs 调整terminal的尺寸
func ResizeWs(c *gin.Context) {
	var req request.ResizeWsStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	clients.lock.RLock()
	cli, ok := clients.data[req.SshId]
	clients.lock.RUnlock()
	if !ok {
		response.FailWithMsg("无法找到连接的ssh实例")
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
		response.FailWithMsg("调整terminal尺寸失败")
	}
}

// GetSshDirInfo 获取文件夹路径下的所有文件信息
func GetSshDirInfo(c *gin.Context) {
	var req request.NodeShellFileStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	clients.lock.RLock()
	cli, ok := clients.data[req.SshId]
	clients.lock.RUnlock()
	if !ok {
		response.FailWithMsg("无法找到连接的ssh实例")
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

	// 内部方法,处理路径信息
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

// GetSshFile 读取文件内容
func GetSshFile(c *gin.Context) {
	var req request.NodeShellFileStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}
	clients.lock.RLock()
	cli, ok := clients.data[req.SshId]
	clients.lock.RUnlock()
	if !ok {
		response.FailWithMsg("无法找到连接的ssh实例")
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

func DownloadFile(c *gin.Context) {
	var req request.NodeShellFileStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	clients.lock.RLock()
	cli, ok := clients.data[req.SshId]
	clients.lock.RUnlock()
	if !ok {
		response.FailWithMsg("无法找到连接的ssh实例")
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

func UpdateFile(c *gin.Context) {
	var req request.ModifyFileRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	clients.lock.RLock()
	cli, ok := clients.data[req.SshId]
	clients.lock.RUnlock()
	if !ok {
		response.FailWithMsg("无法找到连接的ssh实例")
		return
	}
	// 先判断文件是否有读写权限
	file, err := cli.sftpClient.OpenFile(req.Path, os.O_RDWR|os.O_TRUNC)
	if err != nil {
		response.FailWithMsg(fmt.Sprintf("文件%s没有读写权限", req.Path))
		return
	}
	defer file.Close()
	_, err = file.WriteAt([]byte(req.Content), 0)
	if err != nil {
		response.FailWithMsg(fmt.Sprintf("文件%s内容写入失败", file.Name()))
		return
	}
	response.SuccessWithData(fmt.Sprintf("文件%s更新成功", file.Name()))
}

func NodeVncWs(c *gin.Context) {
	var req request.NodeVncWsRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}
	p := vncproxy.New(fmt.Sprintf("%s:%d", req.Address, req.Port))
	handler := nwebsocket.Handler(p.ServeWS)
	handler.ServeHTTP(c.Writer, c.Request)
}

func formatLine(line string) string {
	// Create colorized instances for different styles
	symbolColor := color.New(color.FgGreen)
	monthColor := color.New(color.FgGreen)

	// Define regular expressions for identifying integers, file permissions, and timestamps
	symbolPattern := `-`
	monthPattern := `\b(?:Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)\b`

	// Apply color to file symbol
	newLine := regexp.MustCompile(symbolPattern).ReplaceAllStringFunc(line, func(match string) string {
		return symbolColor.Sprintf("%s", match)
	})

	// Apply color to months
	newLine = regexp.MustCompile(monthPattern).ReplaceAllStringFunc(newLine, func(match string) string {
		return monthColor.Sprintf("%s", match)
	})

	return newLine
}
