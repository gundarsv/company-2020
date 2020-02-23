package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/gundarsv/company-api/handler"
	"github.com/gundarsv/company-api/helpers"
	"github.com/gundarsv/company-api/model"
	"github.com/gundarsv/company-api/repository"
	"gopkg.in/go-playground/validator.v9"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var companies = repository.GetCompanyData()
var owners = repository.GetOwnerData()
var validate *validator.Validate

// GET api/company
func getAllCompanies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(repository.GetAllCompanies())
}

//GET api/owner
func getAllOwners(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(repository.GetAllOwners())
}

// GET api/company/{id}
func getCompanyByID(w http.ResponseWriter, r *http.Request) {
	companyId, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		handler.HandleUserError(w, "Error", 400, err)
		return
	}

	company := repository.GetCompanyByID(companyId)

	if company == nil {
		handler.HandleUserError(w, "Company was not found", 404, errors.New("no company found"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(company)
}

func getOwnerByID(w http.ResponseWriter, r *http.Request) {
	ownerId, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		handler.HandleUserError(w, "Error", 400, err)
		return
	}

	owner := repository.GetOwnerByID(ownerId)

	if owner == nil {
		handler.HandleUserError(w, "Company was not found", 404, errors.New("no company found"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(owner)
}

//POST api/company
func createCompany(w http.ResponseWriter, r *http.Request) {
	reqBody, readErrors := ioutil.ReadAll(r.Body)

	if readErrors != nil {
		handler.HandleUserError(w, "Error", 400, readErrors)
		return
	}

	var newCompany model.Company

	json.Unmarshal(reqBody, &newCompany)

	if validationErrors := validate.Struct(newCompany); validationErrors != nil {
		handler.HandleUserError(w, "Error", 422, validationErrors)
		return
	}

	if newCompany.Owners != nil {
		handler.HandleUserError(w, "Please remove owners", 422, errors.New("owners were added to company"))
		return
	}

	newCompany.ID = repository.GetIdForCompany()
	companies = append(companies, newCompany)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newCompany)
}

//POST api/company/{id}/owner
func addOwnerToCompany(w http.ResponseWriter, r *http.Request) {
	companyId, err := strconv.Atoi(mux.Vars(r)["companyId"])

	if err != nil {
		handler.HandleUserError(w, "Error", 400, err)
		return
	}

	company := helpers.GetCompanyById(companyId, &companies)

	if company == nil {
		handler.HandleUserError(w, "Company was not found", 404, errors.New("no company found"))
		return
	}

	ownerId, err := strconv.Atoi(mux.Vars(r)["companyId"])

	if err != nil {
		handler.HandleUserError(w, "Error", 400, err)
		return
	}

	owner := helpers.GetOwnerById(ownerId, &owners)

	if owner == nil {
		handler.HandleUserError(w, "Owner was not found", 404, errors.New("no owner found"))
		return
	}

	company.AddOwner(*owner)

	repository.UpdateCompany(company)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(company)
}

//func updateCompany(w http.ResponseWriter, r *http.Request) {
//	reqBody, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		log.Println(err)
//		w.WriteHeader(http.StatusBadRequest)
//	}
//
//	var updatedCompany Model.Company
//
//	json.Unmarshal(reqBody, &updatedCompany)
//
//	var company = Helper.GetCompanyById(updatedCompany.ID, companies)
//
//	if company == nil {
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}
//
//	log.Println(updatedCompany)
//
//	if !Helper.AreOwnersInCompanyValid(updatedCompany) {
//		w.WriteHeader(http.StatusUnprocessableEntity)
//		return
//	}
//
//	json.NewEncoder(w).Encode(Helper.UpdateCompany(company, &updatedCompany))
//}



func main() {
	validate = validator.New()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/company", getAllCompanies).Methods("GET")
	router.HandleFunc("/api/company/{id}", getCompanyByID).Methods("GET")
	router.HandleFunc("/api/company", createCompany).Methods("POST")
	//router.HandleFunc("/api/company/{id}", updateCompany).Methods("PUT")
	//router.HandleFunc("/api/company/{id}/owner", updateCompany).Methods("PUT")

	router.HandleFunc("/api/company/{companyId}/owner/{ownerId}", addOwnerToCompany).Methods("PUT")

	router.HandleFunc("/api/owner", getAllOwners).Methods("GET")
	router.HandleFunc("/api/owner/{id}", getOwnerByID).Methods("GET")

	router.Use(LoggingMiddleware)
	log.Println("Server now listening at :8080")
	repository.ConnectToDb()
	log.Fatal(http.ListenAndServe(":8080", router))
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI, r.Method)
		next.ServeHTTP(w, r)
	})
}