package main

import (
	"reflect"
	"fmt"
)

func LoginHandler(client *WsConn, c *Command){
	m := reflect.ValueOf(c.Data).MapIndex(reflect.ValueOf("mobile"))
	otp := reflect.ValueOf(c.Data).MapIndex(reflect.ValueOf("otp"))
	if Password(fmt.Sprintf("%s", m)) == fmt.Sprintf("%s", otp) {
		client.Write(Command{
			Name: "login",
			Data: Data{
				Status: true,
				Message: "Login Successful",
			},
		})
		client.Register()
		client.Authenticate(fmt.Sprintf("%s", m))
		return
	}
	client.Write(Command{
		Name: "login",
		Data: Data{
			Status: false,
			Message: "Invalid Username/Password Combination",
		},
	})
}