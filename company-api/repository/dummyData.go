package repository

import (
	"../model"
	"log"
)

var countForOwners = 1
var countForCompanies = 1

var owners = []model.Owner{
	{
		ID:        GetIdForOwner(),
		FirstName: "Gundars",
		LastName:  "Vasilevskis",
		Address:   "Nylandsalle 50, 8700 Horsens",
	},
	{
		ID:        GetIdForOwner(),
		FirstName: "Gunca",
		LastName:  "Vasilevskis",
		Address:   "Nylandsalle 50, 8700 Horsens",
	},
}

var companies = []model.Company{
	{
		ID:      GetIdForCompany(),
		Name:    "Clearhaus",
		Address: "Navnejausmas 50",
		City:    "Aarhus",
		Country: "Denmark",
	},
	{
		ID:      GetIdForCompany(),
		Name:    "Plan2Learn",
		Address: "Navnejausmas 51",
		City:    "Viby",
		Country: "Denmark",
	},
}

func GetCompanyData() []model.Company {
	companies[1].AddOwner(owners[1])
	return companies
}

func GetOwnerData() []model.Owner {
	return owners
}

func GetIdForCompany() int {
	id := countForCompanies
	countForCompanies++
	return id
}

func GetIdForOwner() int {
	id := countForOwners
	countForOwners++
	return id
}

func UpdateCompany(updatedCompany *model.Company) {
	for _, companyByID := range companies {
		if companyByID.ID == updatedCompany.ID {
			companyByID = *updatedCompany
		}
	}

	log.Println(companies)
}
