package main

import (
	"github.com/gorilla/websocket"
)

type Command struct {
	Name string         `json:"name"`
	Data interface{}    `json:"data"`
}

type Data struct {
	Status bool `json:"status"`
	Message string `json:"message"`
}

func SendError(c *websocket.Conn, e error) {
	c.WriteJSON(Command{
		Name: "error",
		Data: e,
	})
}
