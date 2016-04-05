package main

import (
	"fmt"
	"errors"
)

func VoteHandler(client *WsConn, c *Command){
	if !client.Authenticated || !client.Registered || client.mobile == ""{
		client.Error(errors.New("Unauthorized"))
		return
	}
	currentEvent, _ := redisClient.Get("current").Result()
	if currentEvent == "waiting" {
		client.Error(errors.New("Waiting for next Event"))
		return
	}
	voted, _ := redisClient.Get(client.mobile + ":" + currentEvent).Result()
	if voted != "" {
		client.Error(errors.New("Already Voted"))
		return
	}
	if c.Data != "gs" && c.Data != "lr" {
		client.Error(errors.New("Improper vote"))
		return
	}
	vote := fmt.Sprintf("%s", c.Data)
	p := redisClient.Pipeline()
	set := p.Set(client.mobile + ":" + currentEvent, c.Data, 0)
	counter := p.HIncrBy(currentEvent, "counter", 1)
	voteCounter := p.HIncrBy(currentEvent, vote + "Counter", 1)
	_, err := p.Exec()
	if err != nil || set.Err() != nil || voteCounter.Err() != nil || counter.Err() != nil {
		client.Error(errors.New("Pipeline Error"))
		return
	}
	client.Write(Command{
		Name: "vote",
		Data: "success",
	})
}