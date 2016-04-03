package main

import (
	"fmt"
)

func InitInteractor(connections chan *WsConn) {
	go func() {
		for {
			ws := <- connections
			clients.Add(ws)
			go Launch(ws)
		}
	}()
}

func Launch(ws *WsConn) {
	for {
		var c Command
		if err := ws.Conn.ReadJSON(&c); err != nil {
			ws.Conn.WriteJSON(Command{
				Name: "error",
				Data: err.Error(),
			})
		}
		if c.Name == "echo" {
			fmt.Println("ECHO EVENT")
			ws.Conn.WriteJSON(Command{
				Name: "echo",
				Data: c.Data,
			})
		} else if c.Name == "meta" {
			go MetaHandler(ws.Conn, nil)
		} else if c.Name == "exit" {
			fmt.Println("CLOSE EVENT")
			clients.Remove(ws)
			ws.Conn.Close()
		} else if c.Name == "register" {
			go RegisterHandler(ws, &c)
		} else if c.Name == "login" {
			go LoginHandler(ws, &c)
		} else if c.Name == "event" {
		//	respond with current event
		} else if c.Name == "vote" {
		//	increment counter
		}
	}
}