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
	result, err := redisClient.Set("eventCounter", 1, 0).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
	for _, datum := range strings.Split(stringData, "\n\n") {
		if strings.HasPrefix(datum, "#") {
			continue
		}
		field := strings.Split(datum, "\n")
		if field[0] == "events" {
			for count:=0; count < len(field[1:]); count += 3{
				_, err = redisClient.RPush("events", field[count + 1]).Result()
				if err != nil {
					log.Fatal(err)
				}
				_, err = redisClient.Set(field[count + 1], "waiting", 0).Result()
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