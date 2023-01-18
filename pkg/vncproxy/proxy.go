package vncproxy

import (
	"golang.org/x/net/websocket"
	"metalflow/pkg/global"
)

type Proxy struct {
	Address string
}

func New(addr string) *Proxy {
	return &Proxy{
		Address: addr,
	}
}

func (p *Proxy) ServeWS(ws *websocket.Conn) {
	ws.PayloadType = websocket.BinaryFrame

	peer, err := NewPeer(ws, p.Address)
	if err != nil {
		global.Log.Errorf("get vnc server failed: %v", err)
		return
	}
	defer func(p *Peer) {
		p.Close()
	}(peer)

	go func(p *Peer) {
		e := p.ReadTarget()
		if e != nil {
			global.Log.Errorf("read target stream error: %v", e)
			return
		}
	}(peer)

	err = peer.ReadSource()
	if err != nil {
		global.Log.Errorf("read source stream failed. err: %v", err)
		return
	}
}
