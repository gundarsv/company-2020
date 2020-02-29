package helper

import (
	"encoding/json"
	"log"
	"net/http"
)

type UserMessage struct {
	Message    string
	IsError    bool
	StatusCode int
	Payload    interface{}
	LogMessage string
}

type DatabaseResponse struct {
	IsError     bool
	Message     string
	LogMessage  string
	MessageCode int
}

type ErrorMessage struct {
	Error string
}

type SuccessMessage struct {
	Message string
}

var (
	SomethingWentWrongMessage = UserMessage{
		Message:    "Something went wrong",
		IsError:    true,
		StatusCode: 422,
		Payload:    nil,
	}
)

func NewDatabaseResponse(isError bool, message string, logMessage string, messageCode int) *DatabaseResponse {
	dr := new(DatabaseResponse)
	dr.Message = message
	dr.IsError = isError
	dr.LogMessage = logMessage
	dr.MessageCode = messageCode
	return dr
}

func HandleUserMessage(w http.ResponseWriter, message UserMessage) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(message.StatusCode)

	if message.IsError {
		err := json.NewEncoder(w).Encode(ErrorMessage{Error: message.Message})

		if err != nil {
			log.Println(err.Error())
		}

		log.Println(message.LogMessage)
		return
	}

	if message.Payload == nil {
		err := json.NewEncoder(w).Encode(SuccessMessage{Message: message.Message})

		if err != nil {
			log.Println(err.Error())
		}

		log.Println(message.LogMessage)
		return
	}

	err := json.NewEncoder(w).Encode(message.Payload)

	if err != nil {
		log.Println(err.Error())
	}

	log.Println(message.LogMessage)
	return
}
