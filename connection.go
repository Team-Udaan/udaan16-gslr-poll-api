package main

import (
	"github.com/gorilla/websocket"
	"gopkg.in/mgo.v2/bson"
)

type WsConn struct {
	done            chan bool
	Hash            string
	Conn            *websocket.Conn
	Registered      bool
	Authenticated   bool
}

func NewWsConn(c *websocket.Conn) WsConn {
	hash := bson.NewObjectId().Hex()[0:6]
	return WsConn{
		Hash: hash,
		Conn: c,
		Registered: true,
		Authenticated: false,
	}
}

func (w *WsConn)Register() {
	w.Registered = true
}

func (w *WsConn)Authenticate() {
	w.Authenticated = true
}

func (w *WsConn)Close() error {
	err := w.Conn.Close()
	if err != nil {
		return err
	}
	w.done <- true
	return nil
}

func (w *WsConn)Write(data interface{}) {
	err := w.Conn.WriteJSON(data)
	if err != nil {
		w.Close()
	}
}

type Clients struct {
	Ws map[string]*WsConn
}

func (c *Clients) Add(ws *WsConn) {
	c.Ws[ws.Hash] = ws
	done := make(chan bool)
	ws.done = done
	Launch(ws)
}

func (c *Clients)Remove(ws *WsConn) {
	delete(c.Ws, ws.Hash)
}