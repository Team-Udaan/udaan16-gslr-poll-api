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
	"net/http"
	"net/url"
	"strconv"
	"bytes"
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
func Password(mobile_number string)(string){
	hash := md5.Sum([]byte(mobile_number))
	h := fmt.Sprintf("%x", hash)
	return h[0:6]
}

func message(mobile_number string) (string, error) {
	if utf8.RuneCountInString(mobile_number) != 10 {
		err := errors.New("Invalid Mobile Number/Password")
		return "", err
	}
	return "Thank you for registering! Login with \nUsername:" + mobile_number + "\n Password:" + Password(mobile_number), nil
}

func NewMessage(mobile_number string, message string, custom string, test bool) url.Values {
	v := url.Values{}
	v.Set("hash", TEXTLOCAL_HASH)
	v.Set("username", TEXTLOCAL_USERNAME)
	v.Set("sender", TEXTLOCAL_SENDER)
	v.Set("numbers", mobile_number)
	v.Set("message", message)
	v.Set("test", strconv.FormatBool(test))
	if strings.Compare(custom, "") != 0 {
		v.Set("custom", custom)
	}
	return v
}

type TextlocalError struct {
	Message string `json:"message"`
	Code    int `json:"code"`
}

type TextlocalResponse struct {
	Status string          `json:"status"`
	Errors []TextlocalError  `json:"errors"`
}

func (t *TextlocalResponse)ErrorString() string {
	b := bytes.NewBufferString("")
	encoder := json.NewEncoder(b)
	encoder.Encode(t.Errors)
	return b.String()
}

func (t *TextlocalResponse)String() string {
	b := bytes.NewBufferString("")
	encoder := json.NewEncoder(b)
	encoder.Encode(t)
	return b.String()
}

func SendMessage(m string) (bool, error) {
	msg, err := message(m)
	if err != nil {
		return false, err
	}
	b := NewMessage(m, msg, "", true)
	resp, err := http.PostForm("http://api.textlocal.in/send/", b)
	if err != nil {
		return false, err
	}
	resp_decoder := json.NewDecoder(resp.Body)
	var response TextlocalResponse
	err = resp_decoder.Decode(&response)
	if err != nil {
		return false, err
	} else if response.Status == "failure" {
		return false, errors.New(response.ErrorString())
	}
	return true, nil
}
