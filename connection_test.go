package main

import (
	"testing"
	"github.com/gorilla/websocket"
)

func TestWsConn_Add(t *testing.T) {
	var c websocket.Conn
	var ws = NewWsConn(&c)
	var clients Clients
	clients.Ws = make(map[string]*WsConn)
	clients.Add(&ws)
	if clients.Ws[ws.Hash] == nil {
		t.Fail()
	}

}

func TestWsConn_Remove(t *testing.T) {
	var ws1 WsConn
	var ws2 WsConn
	var ws3 WsConn
	var clients Clients
	clients.Ws = make(map[string]*WsConn)
	ws1.Register()
	ws2.Register()
	ws3.Register()
	clients.Add(&ws1)
	clients.Add(&ws2)
	clients.Add(&ws3)
	clients.Remove(&ws2)
	if clients.Ws[ws2.Hash]!= nil {
		t.Fail()
	}
}