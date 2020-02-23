package model

type Company struct {
	ID          int     `json:"ID"`
	Name        string  `json:"Name" validate:"required"`
	Address     string  `json:"Address" validate:"required"`
	City        string  `json:"City" validate:"required"`
	Country     string  `json:"Country" validate:"required"`
	Email       string  `json:"Email"`
	PhoneNumber string  `json:"PhoneNumber"`
	Owners      []Owner `json:"Owners"`
}

func NewCompany(ID int, Name string, Address string, City string, Country string, Email string, PhoneNumber string) *Company {
	return &Company{ID, Name, Address, City, Country, Email, PhoneNumber, nil}
}

func (c *Company) AddOwner(owner Owner) {
	c.Owners = append(c.Owners, owner)
}
