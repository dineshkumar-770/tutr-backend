package utils

type ResponseStr struct {
	Status   string      `json:"status"`
	Message  string      `json:"message"`
	MyResponse interface{} `json:"response"`
}