package main

import "fmt"

//TODO write a common WriteError function to stop repeating error codes

func VoteHandler(client *WsConn, c *Command){
	if !client.Authenticated || !client.Registered || client.mobile == ""{
		client.Write(Command{
			Name: "error",
			Data: "Unauthorized",
		})
		client.Close()
		return
	}
	currentEvent, _ := redisClient.Get("current").Result()
	if currentEvent == "waiting" {
		client.Write(Command{
			Name: "error",
			Data: "Waiting for next Event",
		})
		client.Close()
		return
	}
	voted, _ := redisClient.Get(client.mobile + ":" + currentEvent).Result()
	if voted != "" {
		client.Write(Command{
			Name: "error",
			Data: "Already Voted",
		})
		return
	}
	if c.Data != "gs" && c.Data != "ls" {
		client.Write(Command{
			Name: "error",
			Data: "Improper Vote",
		})
		client.Close()
		return
	}
	vote := fmt.Sprintf("%s", c.Data)
	p := redisClient.Pipeline()
	set := p.Set(client.mobile + ":" + currentEvent, c.Data, 0)
	hincr := p.HIncrBy(currentEvent, vote + "Votes", 1)
	_, err := p.Exec()
	if err != nil || set.Err() != nil || hincr.Err() != nil {
		client.Write(Command{
			Name: "error",
			Data: err.Error(),
		})
		return
	}
	client.Write(Command{
		Name: "vote",
		Data: "success",
	})
}