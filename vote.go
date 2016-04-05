package main

import (
	"fmt"
	"errors"
)

//TODO write a common WriteError function to stop repeating error codes

func VoteHandler(client *WsConn, c *Command){
	if !client.Authenticated || !client.Registered || client.mobile == ""{
		client.Error(errors.New("Unauthorized"))
		return
	}
	currentEvent, _ := redisClient.Get("current").Result()
	if currentEvent == "waiting" {
		client.Error("Waiting for next Event")
		return
	}
	voted, _ := redisClient.Get(client.mobile + ":" + currentEvent).Result()
	if voted != "" {
		client.Error(errors.New("Already Vote"))
		return
	}
	if c.Data != "gs" && c.Data != "ls" {
		client.Error(errors.New("Improper vote"))
		return
	}
	vote := fmt.Sprintf("%s", c.Data)
	p := redisClient.Pipeline()
	set := p.Set(client.mobile + ":" + currentEvent, c.Data, 0)
	hincr := p.HIncrBy(currentEvent, vote + "Votes", 1)
	_, err := p.Exec()
	if err != nil || set.Err() != nil || hincr.Err() != nil {
		client.Error(errors.New("Pipeline Error"))
		return
	}
	client.Write(Command{
		Name: "vote",
		Data: "success",
	})
}