package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type Error struct {
	Message string
}

func GetError(error string) Error {
	return Error{Message: error}
}

func HandleUserError(w http.ResponseWriter, message string, statusCode int, error error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(GetError(message))
	log.Println(error)
}

func HandleDatabaseError(error error) {
	log.Fatalln(error)
}
