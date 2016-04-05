package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"flag"
	"fmt"
	"os"
	"strings"
	"log"
	"gopkg.in/redis.v3"
)

var clients *Clients = &Clients{}
var Connections chan *WsConn
var CONFIG_FILE *string
var redisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
	Password: "",
	DB: 0,
})
var logger *os.File
var err  error
var panicCount int

func init() {
	logger, err = os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	panicCount = 0
	CONFIG_FILE = flag.String("config", "sample_config.txt", "give file path of configuration file")
	flag.Parse()
	fmt.Println(*CONFIG_FILE)
	Load(*CONFIG_FILE)
	Connections = make(chan *WsConn)
	clients.Ws = make(map[string]*WsConn)
	InitInteractor(Connections)
	TEXTLOCAL_USERNAME = os.Getenv("TEXTLOCAL_USERNAME")
	TEXTLOCAL_HASH = os.Getenv("TEXTLOCAL_HASH")
	TEXTLOCAL_SENDER = os.Getenv("TEXTLOCAL_SENDER")
	if strings.Compare(TEXTLOCAL_USERNAME, "") == 0 && strings.Compare(TEXTLOCAL_HASH, "") == 0 && strings.Compare(TEXTLOCAL_SENDER, "") == 0 {
		log.Fatal("Environment variables not set")
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/open_poll", OpenPollHandler).Methods("POST")
	router.HandleFunc("/ws", WebSocketsHandler)
	http.ListenAndServe(CONFIGURATION["api-ip"][0] + ":" + CONFIGURATION["api-port"][0], router)
}