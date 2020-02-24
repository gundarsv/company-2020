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

func (c *Company) AddOwner(owner Owner) {
	c.Owners = append(c.Owners, owner)
}
