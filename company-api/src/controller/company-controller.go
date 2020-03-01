package controller

import (
	"company-api/src/helper"
	"company-api/src/model"
	"company-api/src/repository"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

func InitCompanyController(router *mux.Router) {
	router.HandleFunc("/api/company", getAllCompanies).Methods(http.MethodGet)
	router.HandleFunc("/api/company/{id}", getCompanyByID).Methods(http.MethodGet)
	router.HandleFunc("/api/company", createCompany).Methods(http.MethodPost)
	router.HandleFunc("/api/company/{id}", updateCompany).Methods(http.MethodPut)
	router.HandleFunc("/api/company/{id}", deleteCompany).Methods(http.MethodDelete)
	router.HandleFunc("/api/company/{companyId}/owner/{ownerId}", addOwnerToCompany).Methods(http.MethodPut)
	router.HandleFunc("/api/company/{companyId}/owner/{ownerId}", deleteOwnerFromCompany).Methods(http.MethodDelete)
}

// GET api/company
func getAllCompanies(w http.ResponseWriter, r *http.Request) {
	dbResponse, companies := repository.GetAllCompanies()

	helper.HandleUserMessage(w, helper.UserMessage{
		Message:    dbResponse.Message,
		IsError:    dbResponse.IsError,
		StatusCode: helper.GetStatusCode(dbResponse.MessageCode),
		Payload:    companies,
		LogMessage: dbResponse.LogMessage,
	})
	return
}

// GET api/company/{id}
func getCompanyByID(w http.ResponseWriter, r *http.Request) {
	companyId, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		helper.SomethingWentWrongMessage.LogMessage = err.Error()
		helper.HandleUserMessage(w, helper.SomethingWentWrongMessage)
		return
	}

	dbResponse, company := repository.GetCompanyByID(companyId)

	helper.HandleUserMessage(w, helper.UserMessage{
		Message:    dbResponse.Message,
		IsError:    dbResponse.IsError,
		StatusCode: helper.GetStatusCode(dbResponse.MessageCode),
		Payload:    company,
		LogMessage: dbResponse.LogMessage,
	})
	return
}

//POST api/company
func createCompany(w http.ResponseWriter, r *http.Request) {
	reqBody, readErrors := ioutil.ReadAll(r.Body)

	if readErrors != nil {
		helper.SomethingWentWrongMessage.LogMessage = readErrors.Error()
		helper.HandleUserMessage(w, helper.SomethingWentWrongMessage)
		return
	}

	var newCompany model.Company

	err := json.Unmarshal(reqBody, &newCompany)

	if err != nil {
		helper.SomethingWentWrongMessage.LogMessage = err.Error()
		helper.HandleUserMessage(w, helper.SomethingWentWrongMessage)
		return
	}

	if validationErrors := validate.Struct(newCompany); validationErrors != nil {
		helper.HandleUserMessage(w, helper.UserMessage{
			Message:    validationErrors.Error(),
			IsError:    true,
			StatusCode: 422,
			Payload:    nil,
			LogMessage: validationErrors.Error(),
		})
		return
	}

	if newCompany.Owners != nil {
		helper.HandleUserMessage(w, helper.UserMessage{
			Message:    "Please remove owners",
			IsError:    true,
			StatusCode: 422,
			Payload:    nil,
			LogMessage: "Owners were added to company",
		})
		return
	}

	dbResponse, createdCompany := repository.CreateCompany(newCompany)

	helper.HandleUserMessage(w, helper.UserMessage{
		Message:    dbResponse.Message,
		IsError:    dbResponse.IsError,
		StatusCode: helper.GetStatusCode(dbResponse.MessageCode),
		Payload:    createdCompany,
		LogMessage: dbResponse.LogMessage,
	})
	return
}

//PUT api/company/{companyID}/owner/{ownerID}
func addOwnerToCompany(w http.ResponseWriter, r *http.Request) {
	companyId, err := strconv.Atoi(mux.Vars(r)["companyId"])

	if err != nil {
		helper.SomethingWentWrongMessage.LogMessage = err.Error()
		helper.HandleUserMessage(w, helper.SomethingWentWrongMessage)
		return
	}

	ownerID, err := strconv.Atoi(mux.Vars(r)["ownerId"])

	if err != nil {
		helper.SomethingWentWrongMessage.LogMessage = err.Error()
		helper.HandleUserMessage(w, helper.SomethingWentWrongMessage)
		return
	}

	dbResponse, company := repository.AddOwnerToCompany(companyId, ownerID)

	helper.HandleUserMessage(w, helper.UserMessage{
		Message:    dbResponse.Message,
		IsError:    dbResponse.IsError,
		StatusCode: helper.GetStatusCode(dbResponse.MessageCode),
		Payload:    company,
		LogMessage: dbResponse.LogMessage,
	})
	return
}

//DELETE api/company/{companyID}/owner/{ownerID}
func deleteOwnerFromCompany(w http.ResponseWriter, r *http.Request) {
	companyId, err := strconv.Atoi(mux.Vars(r)["companyId"])

	if err != nil {
		helper.SomethingWentWrongMessage.LogMessage = err.Error()
		helper.HandleUserMessage(w, helper.SomethingWentWrongMessage)
		return
	}

	ownerID, err := strconv.Atoi(mux.Vars(r)["ownerId"])

	if err != nil {
		helper.SomethingWentWrongMessage.LogMessage = err.Error()
		helper.HandleUserMessage(w, helper.SomethingWentWrongMessage)
		return
	}

	dbResponse := repository.DeleteOwnerFromCompany(companyId, ownerID)

	helper.HandleUserMessage(w, helper.UserMessage{
		Message:    dbResponse.Message,
		IsError:    dbResponse.IsError,
		StatusCode: helper.GetStatusCode(dbResponse.MessageCode),
		Payload:    nil,
		LogMessage: dbResponse.LogMessage,
	})
	return
}

//UPDATE api/company/{id}
func updateCompany(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helper.SomethingWentWrongMessage.LogMessage = err.Error()
		helper.HandleUserMessage(w, helper.SomethingWentWrongMessage)
		return
	}

	var updatedCompany model.Company

	updatedCompany.Name = helper.NilString
	updatedCompany.Address = helper.NilString
	updatedCompany.Country = helper.NilString
	updatedCompany.City = helper.NilString
	updatedCompany.Email = helper.NilString
	updatedCompany.PhoneNumber = helper.NilString

	if err := json.Unmarshal(reqBody, &updatedCompany); err != nil {
		helper.SomethingWentWrongMessage.LogMessage = err.Error()
		helper.HandleUserMessage(w, helper.SomethingWentWrongMessage)
		return
	}

	dbResponse, company := repository.UpdateCompany(id, updatedCompany)

	helper.HandleUserMessage(w, helper.UserMessage{
		Message:    dbResponse.Message,
		IsError:    dbResponse.IsError,
		StatusCode: helper.GetStatusCode(dbResponse.MessageCode),
		Payload:    company,
		LogMessage: dbResponse.LogMessage,
	})
	return
}

//DELETE api/company/{id}
func deleteCompany(w http.ResponseWriter, r *http.Request) {
	companyId, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		helper.SomethingWentWrongMessage.LogMessage = err.Error()
		helper.HandleUserMessage(w, helper.SomethingWentWrongMessage)
		return
	}

	dbResponse := repository.DeleteCompany(companyId)

	helper.HandleUserMessage(w, helper.UserMessage{
		Message:    dbResponse.Message,
		IsError:    dbResponse.IsError,
		StatusCode: helper.GetStatusCode(dbResponse.MessageCode),
		Payload:    nil,
		LogMessage: dbResponse.LogMessage,
	})
	return
}
