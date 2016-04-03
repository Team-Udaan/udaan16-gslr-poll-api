package main

import (
	"net/http"
	"os"
	"encoding/json"
	"fmt"
)

func OpenPollHandler(w http.ResponseWriter, r *http.Request) {
	username := os.Getenv("username")
	password := os.Getenv("password")
	reqBodyDecoder := json.NewDecoder(r.Body)
	var reqBody map[string]string
	if err := reqBodyDecoder.Decode(&reqBody); err != nil {
		RespondHTTP(w, err.Error(), 500)
		return
	}
	fmt.Println(reqBody)
	if(username != reqBody["username"] ||
		password != reqBody["password"]) {
		RespondHTTP(w, "Unauthorized", 401)
		return
	}
	if reqBody["event"] == "" {
		RespondHTTP(w, "'event' parameter missing in request body", 400)
		return
	}
	event, err := redisClient.Get(reqBody["event"]).Result()
	if event == "" {
		RespondHTTP(w, "No such event registered", 400)
		return
	}
	if err != nil {
		RespondHTTP(w, err.Error(), 500)
		return
	}
	result, err := redisClient.Set(event, "polling", 0).Result()
	if err != nil {
		RespondHTTP(w, err.Error(), 500)
		return
	}
	fmt.Println("Result ", result)
	current, err := redisClient.Get("current").Result()
	if err != nil {
		RespondHTTP(w, err.Error(), 500)
		return
	}
	if current != "waiting" {
		clients.Broadcast(Command{
			Name: "event",
			Data: Command{
				Name: current,
			Data: "close",
			},
		})
	}
	result, err = redisClient.Set("current", reqBody["event"], 0).Result()
	if err != nil {
		RespondHTTP(w, err.Error(), 500)
		return
	}
	// TODO Broadcast the opening of current event
	clients.Broadcast(Command{
			Name: "event",
			Data: Command{
				Name: reqBody["event"],
			Data: "open",
			},
		})
	fmt.Println("Result ", result)
	RespondHTTP(w, reqBody["event"] + " open for polling", 200)
}
