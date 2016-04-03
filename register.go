package main

import (
	"reflect"
	"fmt"
)

func RegisterHandler(client *WsConn, c *Command) {
	m := reflect.ValueOf(c.Data).MapIndex(reflect.ValueOf("mobile"))
	if ok, err := SendMessage(fmt.Sprintf("%s", m)); !ok {
		fmt.Println(err)
	}
	client.Registered = true
	return
}