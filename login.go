package main

import (
	"reflect"
	"fmt"
)

func LoginHandler(client *WsConn, c *Command){
	m := reflect.ValueOf(c.Data).MapIndex(reflect.ValueOf("mobile"))
	otp := reflect.ValueOf(c.Data).MapIndex(reflect.ValueOf("otp"))
	if Password(fmt.Sprintf("%s", m)) == fmt.Sprintf("%s", otp) {
		err := client.Conn.WriteJSON(Command{
			Name: "login",
			Data: Data{
				Status: true,
				Message: "Login Successful",
			},
		})
		if err != nil {
			fmt.Println(err)
			client.Conn.Close()
			clients.Remove(client)
			return
		}
		client.Register()
		client.Authenticate()
		return
	}
	err := client.Conn.WriteJSON(Command{
		Name: "login",
		Data: Data{
			Status: false,
			Message: "Invalid Username/Password Combination",
		},
	})
	if err != nil {
		fmt.Println(err)
	}
	client.Conn.Close()
	clients.Remove(client)
}