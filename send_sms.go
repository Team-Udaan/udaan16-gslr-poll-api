package main

import (
	"errors"
	"crypto/md5"
	"fmt"
	"strings"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"bytes"
)

var TEXTLOCAL_USERNAME string
var TEXTLOCAL_HASH string
var TEXTLOCAL_SENDER string

func Password(mobile_number string)(string){
	hash := md5.Sum([]byte(mobile_number))
	h := fmt.Sprintf("%x", hash)
	return h[0:4]
}

func message(mobileNumber string) (string, error) {
	if n := len(mobileNumber); n != 10 {
		err := errors.New("Invalid Mobile Number " + strconv.Itoa(n))
		return "", err
	}
	return "Thank you for registering in GS vs LR Live Online Poll! Login with \nUsername: " + mobileNumber + "\nPassword: " + Password(mobileNumber), nil
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
	//TODO Remove true in production
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
