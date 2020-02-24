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

func InitCompanyController(router *mux.Router) {
	router.HandleFunc("/api/company", getAllCompanies).Methods("GET")
	router.HandleFunc("/api/company/{id}", getCompanyByID).Methods("GET")
	//router.HandleFunc("/api/company/{id}", updateCompany).Methods("PUT")
	//router.HandleFunc("/api/company", createCompany).Methods("POST")
	//router.HandleFunc("/api/company/{companyId}/owner/{ownerId}", addOwnerToCompany).Methods("PUT")
}

// GET api/company
func getAllCompanies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(repository.GetAllCompanies())
}

// GET api/company/{id}
func getCompanyByID(w http.ResponseWriter, r *http.Request) {
	companyId, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		helper.HandleUserError(w, "Error", 400, err)
		return
	}

	company := repository.GetCompanyByID(companyId)

	if company == nil {
		helper.HandleUserError(w, "Company was not found", 404, errors.New("no company found"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(company)
}

////POST api/company
//func createCompany(w http.ResponseWriter, r *http.Request) {
//	reqBody, readErrors := ioutil.ReadAll(r.Body)
//
//	if readErrors != nil {
//		helper.HandleUserError(w, "Error", 400, readErrors)
//		return
//	}
//
//	var newCompany model.Company
//
//	json.Unmarshal(reqBody, &newCompany)
//
//	if validationErrors := validate.Struct(newCompany); validationErrors != nil {
//		helper.HandleUserError(w, "Error", 422, validationErrors)
//		return
//	}
//
//	if newCompany.Owners != nil {
//		helper.HandleUserError(w, "Please remove owners", 422, errors.New("owners were added to company"))
//		return
//	}
//
//	newCompany.ID = repository.GetIdForCompany()
//	companies = append(companies, newCompany)
//
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(newCompany)
//}

//POST api/company/{id}/owner
//func addOwnerToCompany(w http.ResponseWriter, r *http.Request) {
//	companyId, err := strconv.Atoi(mux.Vars(r)["companyId"])
//
//	if err != nil {
//		helper.HandleUserError(w, "Error", 400, err)
//		return
//	}
//
//	company := helper.GetCompanyById(companyId, &companies)
//
//	if company == nil {
//		helper.HandleUserError(w, "Company was not found", 404, errors.New("no company found"))
//		return
//	}
//
//	ownerId, err := strconv.Atoi(mux.Vars(r)["companyId"])
//
//	if err != nil {
//		helper.HandleUserError(w, "Error", 400, err)
//		return
//	}
//
//	owner := helper.GetOwnerById(ownerId, &owners)
//
//	if owner == nil {
//		helper.HandleUserError(w, "Owner was not found", 404, errors.New("no owner found"))
//		return
//	}
//
//	company.AddOwner(*owner)
//
//	repository.UpdateCompany(company)
//
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(company)
//}

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
