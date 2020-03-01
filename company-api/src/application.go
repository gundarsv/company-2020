package main

import (
	"company-api/src/controller"
	"company-api/src/helper"
	"company-api/src/repository"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

func main() {

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"*"},
	})

	router := mux.NewRouter().StrictSlash(true)

	controller.InitCompanyController(router)
	controller.InitOwnerController(router)

	router.Use(helper.InitLoggingMiddleware)

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/")))
	router.PathPrefix("/company").Handler(http.FileServer(http.Dir("./web/")))
	router.PathPrefix("/owner").Handler(http.FileServer(http.Dir("./web/")))

	log.Println("Server now listening at :" + os.Getenv("PORT"))
	repository.InitRepository()
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), helper.RemoveTrailingSlash(c.Handler(router))))
}
