package main

import (
	"net/http"
	"os"
	"encoding/json"
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
	if(username != reqBody["username"] ||
		password != reqBody["password"]) {
		RespondHTTP(w, "Unauthorized", 401)
		return
	}
	if reqBody["event"] == "" {
		RespondHTTP(w, "'event' parameter missing in request body", 400)
		return
	}
	event, err := redisClient.HGetAll(reqBody["event"]).Result()
	if len(event) == 0 {
		RespondHTTP(w, "No such event registered", 400)
		return
	}
	if err != nil {
		RespondHTTP(w, err.Error(), 500)
		return
	}
	_, err = redisClient.HSet(reqBody["event"], "status", "polling").Result()
	if err != nil {
		RespondHTTP(w, err.Error(), 500)
		return
	}
	_, err = redisClient.Set("current", reqBody["event"], 0).Result()
	if err != nil {
		RespondHTTP(w, err.Error(), 500)
		return
	}
	eventResponse, err := EventResponse(reqBody["event"])
	if err != nil {
		RespondHTTP(w, err.Error(), 500)
		return
	}
	// TODO Broadcast the opening of current event
	clients.Broadcast(Command{
			Name: "event",
			Data: eventResponse,
		})
	RespondHTTP(w, reqBody["event"] + " open for polling", 200)
}
