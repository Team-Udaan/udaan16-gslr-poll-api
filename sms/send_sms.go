package sms

import (
	"unicode/utf8"
	"errors"
	"crypto/md5"
	"fmt"
	"os"
	"strings"
	"log"
	"encoding/json"
	"github.com/parnurzeal/gorequest"
)

var TEXTLOCAL_USERNAME string
var TEXTLOCAL_HASH string
var TEXTLOCAL_SENDER string

func init()  {
	TEXTLOCAL_USERNAME = os.Getenv("TEXTLOCAL_USERNAME")
	TEXTLOCAL_HASH = os.Getenv("TEXTLOCAL_HASH")
	TEXTLOCAL_SENDER = os.Getenv("TEXTLOCAL_SENDER")
	if strings.Compare(TEXTLOCAL_USERNAME, "") == 0 && strings.Compare(TEXTLOCAL_HASH, "") == 0 && strings.Compare(TEXTLOCAL_SENDER, "") == 0 {
		log.Fatal("Environment variables not set")
	}
}
func GetPassword(mobile_number string)(string){
	hash := md5.Sum([]byte(mobile_number))
	h := fmt.Sprintf("%x", hash)
	return h[0:6]
}

func getMessage(mobile_number string) (string, error) {
	if utf8.RuneCountInString(mobile_number) != 10 {
		err := errors.New("Invalid Mobile Number/Password")
		return "", err
	}
	return "Thank you for registering! Login with \nUsername:" + mobile_number + "\n Password:" + GetPassword(mobile_number), nil
}

type Message struct {
	Hash     string  `json:"hash"`
	Sender   string  `json:"sender"`
	Username string  `json:"username"`
	Text     string  `json:"message"`
	Custom   string  `json:"custom,omitempty"`
	Test     bool    `json:"test,omitempty"`
}

func NewMessage(message string, custom string, test bool) Message {
	return Message{
		Hash: TEXTLOCAL_HASH,
		Sender: TEXTLOCAL_SENDER,
		Username: TEXTLOCAL_USERNAME,
		Text: message,
		Custom: custom,
		Test: test,
	}
}

func SendMessage(m string) (bool, error) {
	mes, err := getMessage(m)
	if err != nil {
		return false, err
	}
	b := NewMessage(mes, "", true)
	c := gorequest.New()
	r := c.Post("http://api.textlocal.in/send/")
	body, err := json.Marshal(&b)
	os.Stdout.Write(body)
	if err != nil {
		return false, err
	}
	req_body := string(body[0:len(body)])
	r.Send(req_body)
	_, resp_body, errs := r.End()
	fmt.Println(resp_body, errs)
	return true, nil
}
