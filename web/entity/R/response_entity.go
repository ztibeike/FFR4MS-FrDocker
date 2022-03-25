package R

import "net/http"

type ResponseEntity struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func OK(data interface{}) *ResponseEntity {
	entity := &ResponseEntity{
		Code:    http.StatusOK,
		Message: "Success!",
		Data:    data,
	}
	return entity
}

func Error(code int, message string, data interface{}) *ResponseEntity {
	if message == "" {
		message = http.StatusText(code)
	}
	entity := &ResponseEntity{
		Code:    code,
		Message: message,
		Data:    data,
	}
	return entity
}
