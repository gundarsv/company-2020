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

func InitOwnerController(router *mux.Router) {
	router.HandleFunc("/api/owner", getAllOwners).Methods(http.MethodGet)
	router.HandleFunc("/api/owner/{id}", getOwnerByID).Methods(http.MethodGet)
	router.HandleFunc("/api/owner", createOwner).Methods(http.MethodPost)
	router.HandleFunc("/api/owner/{id}", updateOwner).Methods(http.MethodPut)
	router.HandleFunc("/api/owner/{id}", deleteOwner).Methods(http.MethodDelete)
}

//GET api/owner
func getAllOwners(w http.ResponseWriter, r *http.Request) {
	dbResponse, owners := repository.GetAllOwners()

	helper.HandleUserMessage(w, helper.UserMessage{
		Message:    dbResponse.Message,
		IsError:    dbResponse.IsError,
		StatusCode: helper.GetStatusCode(dbResponse.MessageCode),
		Payload:    owners,
		LogMessage: dbResponse.LogMessage,
	})
	return
}

//GET api/owner/{id}
func getOwnerByID(w http.ResponseWriter, r *http.Request) {
	ownerId, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		helper.SomethingWentWrongMessage.LogMessage = err.Error()
		helper.HandleUserMessage(w, helper.SomethingWentWrongMessage)
		return
	}

	dbResponse, owner := repository.GetOwnerByID(ownerId)

	helper.HandleUserMessage(w, helper.UserMessage{
		Message:    dbResponse.Message,
		IsError:    dbResponse.IsError,
		StatusCode: helper.GetStatusCode(dbResponse.MessageCode),
		Payload:    owner,
		LogMessage: dbResponse.LogMessage,
	})
	return
}

//POST api/owner
func createOwner(w http.ResponseWriter, r *http.Request) {
	reqBody, readErrors := ioutil.ReadAll(r.Body)

	if readErrors != nil {
		helper.SomethingWentWrongMessage.LogMessage = readErrors.Error()
		helper.HandleUserMessage(w, helper.SomethingWentWrongMessage)
		return
	}

	var newOwner model.Owner

	err := json.Unmarshal(reqBody, &newOwner)

	if err != nil {
		helper.SomethingWentWrongMessage.LogMessage = err.Error()
		helper.HandleUserMessage(w, helper.SomethingWentWrongMessage)
		return
	}

	if validationErrors := validate.Struct(newOwner); validationErrors != nil {
		helper.HandleUserMessage(w, helper.UserMessage{
			Message:    validationErrors.Error(),
			IsError:    true,
			StatusCode: 422,
			Payload:    nil,
			LogMessage: validationErrors.Error(),
		})
		return
	}

	dbResponse, createdOwner := repository.CreateOwner(newOwner)

	helper.HandleUserMessage(w, helper.UserMessage{
		Message:    dbResponse.Message,
		IsError:    dbResponse.IsError,
		StatusCode: helper.GetStatusCode(dbResponse.MessageCode),
		Payload:    createdOwner,
		LogMessage: dbResponse.LogMessage,
	})
	return
}

//UPDATE api/owner/{id}
func updateOwner(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helper.SomethingWentWrongMessage.LogMessage = err.Error()
		helper.HandleUserMessage(w, helper.SomethingWentWrongMessage)
		return
	}

	var updatedOwner model.Owner

	updatedOwner.FirstName = helper.NilString
	updatedOwner.LastName = helper.NilString
	updatedOwner.Address = helper.NilString

	if err := json.Unmarshal(reqBody, &updatedOwner); err != nil {
		helper.SomethingWentWrongMessage.LogMessage = err.Error()
		helper.HandleUserMessage(w, helper.SomethingWentWrongMessage)
		return
	}

	dbResponse, owner := repository.UpdateOwner(id, updatedOwner)

	helper.HandleUserMessage(w, helper.UserMessage{
		Message:    dbResponse.Message,
		IsError:    dbResponse.IsError,
		StatusCode: helper.GetStatusCode(dbResponse.MessageCode),
		Payload:    owner,
		LogMessage: dbResponse.LogMessage,
	})
	return
}

//DELETE api/owner/{id}
func deleteOwner(w http.ResponseWriter, r *http.Request) {
	companyId, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		helper.SomethingWentWrongMessage.LogMessage = err.Error()
		helper.HandleUserMessage(w, helper.SomethingWentWrongMessage)
		return
	}

	dbResponse := repository.DeleteOwner(companyId)

	helper.HandleUserMessage(w, helper.UserMessage{
		Message:    dbResponse.Message,
		IsError:    dbResponse.IsError,
		StatusCode: helper.GetStatusCode(dbResponse.MessageCode),
		Payload:    nil,
		LogMessage: dbResponse.LogMessage,
	})
	return

}
