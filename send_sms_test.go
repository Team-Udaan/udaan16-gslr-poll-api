package main

import (
	"testing"
	"strings"
)

func TestGetMessage(t *testing.T)  {
	mobile_number := "7600647682"
	expected := "Thank you for registering! Login with \nUsername:" + mobile_number + "\n Password:" + Password(mobile_number)
	actual, err := message(mobile_number)
	if err == nil && strings.Compare(expected, actual) != 0 {
		t.Fail()
	}
	mobile_number = "9089789"
	actual, err = message(mobile_number)
	expected = ""
	if strings.Compare(err.Error(), "Invalid Movile Number/Password") != 0  && strings.Compare(expected, actual) != 0 {
		t.Fail()
	}
}
