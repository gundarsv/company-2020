package main

import (
	"company-api/src/controller"
	"company-api/src/helper"
	"company-api/src/repository"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Not found test")
}

func main() {

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"*"},
	})

	router := mux.NewRouter().StrictSlash(true)

	controller.InitCompanyController(router)
	controller.InitOwnerController(router)

	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})

	router.Use(helper.InitLoggingMiddleware)

	log.Println("Server now listening at :" + os.Getenv("PORT"))
	repository.InitRepository()
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), helper.RemoveTrailingSlash(c.Handler(router))))
}
