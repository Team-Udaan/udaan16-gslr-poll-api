package handlers

import (
	"encoding/json"
	"net/http"
	"fmt"
)

type Response struct {
	Message interface{} `json:"message"`
	Code    int         `json:"code"`
}

func Respond(w http.ResponseWriter, response interface{}, code int) {
	//buffer := bytes.NewBufferString(response)
	headers := w.Header()
	headers.Add("Content-Type", "application/json")
	r := Response{
		Message: response,
		Code: code,
	}
	fmt.Println(r)
	encoder := json.NewEncoder(w)
	err := encoder.Encode(&r)
	if err != nil {
		r.Message = err.Error()
		r.Code = 500
		fmt.Fprint(w, r)
	}
}
