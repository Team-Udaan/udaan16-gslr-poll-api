package main

type Command struct {
	Name string         `json:"name"`
	Data interface{}    `json:"data,omitempty"`
}

type Data struct {
	Status bool `json:"status"`
	Message string `json:"message"`
}
