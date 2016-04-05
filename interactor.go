package main

import (
	"fmt"
	"strconv"
)

func InitInteractor(connections chan *WsConn) {
	go func() {
		for {
			ws := <- connections
			clients.Add(ws)
			//go Launch(ws)
		}
	}()
}

//TODO panic recovery

func Launch(ws *WsConn) {
	defer func() {
		if r := recover(); r != nil {
			logger.WriteString(strconv.Itoa(panicCount)+ ") " +
			fmt.Sprintf("%s", r) + "\n" +
			fmt.Sprintf("%s", ws.mobile) + "\n\n")
			ws.Close()
			clients.Remove(ws)
			panicCount++
		}
	}()
	for {
		select {
		case value := <- ws.done:
			if value == true {
				clients.Remove(ws)
				return
			}
		default:
			var c Command
			if err := ws.Conn.ReadJSON(&c); err != nil {
				ws.Error(err)
			}
			if c.Name == "echo" {
				fmt.Println("ECHO EVENT")
				ws.Write(Command{
					Name: "echo",
					Data: c.Data,
				})
			} else if c.Name == "meta" {
				go MetaHandler(ws, nil)
			} else if c.Name == "exit" {
				ws.Close()
			} else if c.Name == "register" {
				go RegisterHandler(ws, &c)
			} else if c.Name == "login" {
				//Not safe to make this a goroutine
				LoginHandler(ws, &c)
			} else if c.Name == "event" {
				EventHandler(ws, &c)
			} else if c.Name == "vote" {
				VoteHandler(ws, &c)
			}
		}
	}
}