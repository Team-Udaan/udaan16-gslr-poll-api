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
	for _, datum := range strings.Split(stringData, "\n\n") {
		if strings.HasPrefix(datum, "#") {
			continue
		}
		field := strings.Split(datum, "\n")
		config[field[0]] = field[1:]
		//fmt.Println(field)
	}
	fmt.Println(config)
	CONFIGURATION = config
}