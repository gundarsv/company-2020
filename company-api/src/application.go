package main

import (
	"company-api/src/controller"
	"company-api/src/helper"
	"company-api/src/repository"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)

	controller.InitCompanyController(router)
	controller.InitOwnerController(router)

	router.Use(helper.InitLoggingMiddleware)

	log.Println("Server now listening at :" + os.Getenv("PORT"))
	repository.InitRepository()
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), helper.RemoveTrailingSlash(router)))
}
