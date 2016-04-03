package main

type teams struct {
	Name   string `json:"name"`
	Mascot string `json:"mascot"`
}

type Meta struct {
	NightName string `json:"name"`
	Year string `json:"year"`
	GS   teams `json:"gs"`
	LR   teams `json:"lr"`
}

func MetaHandler(conn *WsConn, c *Command)  {
	conn.Write(Command{
		Name: "meta",
		Data: Meta{
			NightName: CONFIGURATION["night-name"][0],
			Year: CONFIGURATION["year"][0],
			GS: teams{
				Name: CONFIGURATION["gs"][0],
				Mascot: CONFIGURATION["gs-mascot"][0],
			},
			LR: teams{
				Name: CONFIGURATION["lr"][0],
				Mascot: CONFIGURATION["lr-mascot"][0],
			},
		},
	})
}