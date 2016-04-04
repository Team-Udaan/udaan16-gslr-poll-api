package main

func VoteHandler(client *WsConn, c *Command){
	if !client.Authenticated || !client.Registered || client.mobile == ""{
		client.Write(Command{
			Name: "error",
			Data: "Unauthorized",
		})
		client.Close()
		return
	}
	voted, _ := redisClient.Get(client.mobile).Result()
	if voted != "" {
		client.Write(Command{
			Name: "error",
			Data: "Already Voted",
		})
		return
	}
	_, err := redisClient.Set(client.mobile, c.Data, 0).Result()
	if err != nil {
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