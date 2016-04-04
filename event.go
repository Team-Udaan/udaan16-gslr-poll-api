package main

type Event struct {
	Name string     `json:"name"`
	GS   Command    `json:"gs"`
	LR   Command    `json:"lr"`
}

func EventResponse(event string) (e *Event, err error) {
	var currentEvent string
	if event == "" {
		currentEvent, err = redisClient.Get("current").Result()
	} else {
		currentEvent = event
	}
	if err != nil{
		return nil, err
	}
	gsTeamName, err := redisClient.HGet(currentEvent, "gs").Result()
	if err != nil{
		return nil, err
	}
	lrTeamName, err := redisClient.HGet(currentEvent, "lr").Result()
	if err != nil{
		return nil, err
	}
	return &Event{
		Name: currentEvent,
		GS: Command{
			Name: gsTeamName,
		},
		LR: Command{
			Name: lrTeamName,
		},
	}, nil
}

func EventHandler(client *WsConn, c *Command) {
	current, err := redisClient.Get("current").Result()
	if err != nil {
		client.Write(Command{
			Name: "error",
			Data: err.Error(),
		})
		return
	}
	voted, _ := redisClient.Get(client.mobile).Result()
	if current == "waiting" || voted != "" {
		client.Write(Command{
			Name: "event",
			Data: "waiting",
		})
		return
	}
	eventResponse, err := EventResponse(current)
	if err != nil {
		client.Write(Command{
			Name: "error",
			Data: err.Error(),
		})
		return
	}
	client.Write(Command{
		Name: "event",
		Data: eventResponse,
	})
}