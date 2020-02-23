package helpers

import (
	"../model"
)

func GetCompanyById(companyID int, companies *[]model.Company) *model.Company {
	for _, companyByID := range *companies {
		if companyByID.ID == companyID {
			return &companyByID
		}
	}

	return nil
}

func GetOwnerById(ownerID int, owners *[]model.Owner) *model.Owner {
	for _, ownerByID := range *owners {
		if ownerByID.ID == ownerID {
			return &ownerByID
		}
	}

	return nil
}
