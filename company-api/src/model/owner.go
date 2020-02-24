package model

type Owner struct {
	ID        int    `json:"ID"`
	FirstName string `json:"FirstName" validate:"required"`
	LastName  string `json:"LastName" validate:"required"`
	Address   string `json:"Address" validate:"required"`
}
