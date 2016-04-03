package main

import (
	"net/http"
	"github.com/gorilla/websocket"
	"fmt"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WebSocketsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}
	ws := NewWsConn(conn)
	clients.Add(&ws)
	Connections <- &ws
}