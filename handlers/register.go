package handlers

import (
	"net/http"
	"encoding/json"
	"github.com/Team-Udaan/udaan16-gslr-poll-api/sms"
)

type RegisterRequestBody struct {
	Mobile string `json:"mobile"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	var reqBody RegisterRequestBody
	err := d.Decode(&reqBody)
	if err != nil {
		Respond(w, err.Error(), 500)
	}
	ok, err := sms.SendMessage(reqBody.Mobile)
	if ok {
		Respond(w, true, 200)
	} else {
		Respond(w, err.Error(), 500)
	}
}

