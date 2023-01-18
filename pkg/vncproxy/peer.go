package vncproxy

import (
	"errors"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"metalflow/pkg/global"
	"net"
	"time"
)

// Peer represents a vnc proxy Peer
// with a websocket connection and a vnc backend connection
type Peer struct {
	source *websocket.Conn
	target net.Conn
}

func NewPeer(ws *websocket.Conn, addr string) (*Peer, error) {
	var (
		c   net.Conn
		err error
	)

	if ws == nil {
		return nil, errors.New("websocket connection is nil")
	}
	c, err = net.DialTimeout("tcp", addr, time.Duration(global.Conf.System.ConnectTimeout)*time.Second)
	if err != nil {
		return nil, fmt.Errorf("can't connect to vnc server, err:%v", err)
	}
	err = c.(*net.TCPConn).SetKeepAlive(true)
	if err != nil {
		return nil, fmt.Errorf("enable vnc server connection keepalive failed, err:%v", err)
	}
	err = c.(*net.TCPConn).SetKeepAlivePeriod(time.Duration(global.Conf.System.ExecuteTimeout) * time.Second)
	if err != nil {
		return nil, fmt.Errorf("set vnc server connetion keepalive period failed, err:%v", err)
	}
	c, err = Connect(addr, ws, c)
	if err != nil {
		return nil, err
	}
	return &Peer{
		source: ws,
		target: c,
	}, nil
}

// ReadSource copy source stream to target connection.
func (p *Peer) ReadSource() error {
	_, err := io.Copy(p.target, p.source)
	if err != nil {
		return fmt.Errorf("copy source(%v) => target(%v) failed", p.source.RemoteAddr(), p.target.RemoteAddr())
	}
	return nil
}

// ReadTarget copy target stream to source connection.
func (p *Peer) ReadTarget() error {
	if _, err := io.Copy(p.source, p.target); err != nil {
		return fmt.Errorf("copy target(%v) => source(%v) failed", p.target.RemoteAddr(), p.source.RemoteAddr())
	}
	return nil
}

// Close the websocket connection and the vnc backend connection.
func (p *Peer) Close() {
	_ = p.source.Close()
	_ = p.target.Close()
}
