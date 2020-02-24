package main

import (
	"company-api/src/controller"
	"company-api/src/helper"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)

	controller.InitCompanyController(router)
	controller.InitOwnerController(router)

	router.Use(helper.InitLoggingMiddleware)

	log.Println("Server now listening at :8080")
	//repository.InitRepository()
	log.Fatal(http.ListenAndServe(":8080", router))
}
