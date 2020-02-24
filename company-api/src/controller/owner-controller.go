package controller

import (
	"company-api/src/helper"
	"company-api/src/repository"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func InitOwnerController(router *mux.Router) {
	router.HandleFunc("/api/owner", getAllOwners).Methods("GET")
	router.HandleFunc("/api/owner/{id}", getOwnerByID).Methods("GET")

}

//GET api/owner
func getAllOwners(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(repository.GetAllOwners())
}

//GET api/owner/{id}
func getOwnerByID(w http.ResponseWriter, r *http.Request) {
	ownerId, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		helper.HandleUserError(w, "Error", 400, err)
		return
	}

	owner := repository.GetOwnerByID(ownerId)

	if owner == nil {
		helper.HandleUserError(w, "Company was not found", 404, errors.New("no company found"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(owner)
}
