package main

import (
	"os"
	"log"
	"io/ioutil"
	"fmt"
	"bytes"
	"strings"
)

var CONFIGURATION map[string][]string

type Loader interface {
	Load(fileName string) map[string]string
}

func Load(fileName string){
	config := make(map[string][]string)
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	byteConfig, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	buffer := bytes.NewBuffer(byteConfig)
	stringData := buffer.String()
	//TODO Add functionality where a "next" keyword in openpoll body, stops the current event and starts the next in line
	//_, err = redisClient.Set("eventCounter", 0, 0).Result()
	//if err != nil {
	//	log.Fatal(err)
	//}
	_, err = redisClient.Set("current", "waiting", 0).Result()
	if err != nil {
		log.Fatal(err)
	}
	for _, datum := range strings.Split(stringData, "\n\n") {
		if strings.HasPrefix(datum, "#") {
			continue
		}
		field := strings.Split(datum, "\n")
		if field[0] == "events" {
			events := field[1:]
			for count:=0; count < len(events); count += 3{
				_, err = redisClient.RPush("events", events[count]).Result()
				if err != nil {
					log.Fatal(err)
				}
				gsTeamName := strings.TrimPrefix(events[count + 1], "GS: ")
				lrTeamName := strings.TrimPrefix(events[count + 2], "LR: ")
				_, err = redisClient.HSet(events[count], "gs", gsTeamName).Result()
				if err != nil {
					log.Fatal(err)
				}
				_, err = redisClient.HSet(events[count], "lr", lrTeamName).Result()
				if err != nil {
					log.Fatal(err)
				}
				_, err = redisClient.HSet(events[count], "status", "waiting").Result()
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		config[field[0]] = field[1:]
	}
	fmt.Println(config)
	CONFIGURATION = config
}