package repository

import (
	"company-api/src/helper"
	"company-api/src/model"
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"os"
)

var (
	server             = os.Getenv("DATABASE_IP")
	port               = os.Getenv("DATABASE_PORT")
	user               = "sa"
	password           = "Secret!Secret"
	database           = "companyDB"
	databaseConnection *sql.DB
)

func GetAllOwners() []*model.Owner {
	rows, err := databaseConnection.Query("SELECT ID, FirstName, LastName, Address FROM dbo.owner;")

	if err != nil {
		helper.HandleDatabaseError(err)
	}
	defer rows.Close()

	var owners []*model.Owner

	for rows.Next() {
		o := new(model.Owner)
		err := rows.Scan(&o.ID, &o.FirstName, &o.LastName, &o.Address)

		if err != nil {
			helper.HandleDatabaseError(err)
		}

		owners = append(owners, o)
	}

	return owners
}

func GetAllCompanies() []*model.Company {
	rows, err := databaseConnection.Query("SELECT ID, Name, Address, City, Country, COALESCE(Email, ''), COALESCE(PhoneNumber, '') FROM dbo.company;")

	if err != nil {
		helper.HandleDatabaseError(err)
	}
	defer rows.Close()

	var companies []*model.Company

	for rows.Next() {
		c := new(model.Company)
		err := rows.Scan(&c.ID, &c.Name, &c.Address, &c.City, &c.Country, &c.PhoneNumber, &c.Email)

		if err == sql.ErrNoRows {
			return nil
		}

		if err != nil {
			helper.HandleDatabaseError(err)
		}

		companies = append(companies, c)
	}

	return companies
}

func GetCompanyByID(companyID int) *model.Company {
	rows := databaseConnection.QueryRow("SELECT ID, Name, Address, City, Country, COALESCE(Email, ''), COALESCE(PhoneNumber, '') FROM dbo.company WHERE ID = ?;", companyID)

	company := new(model.Company)

	if err := rows.Scan(&company.ID, &company.Name, &company.Address, &company.City, &company.Country, &company.Email, &company.PhoneNumber); err != nil {
		if err == sql.ErrNoRows {
			return nil
		} else {
			helper.HandleDatabaseError(err)
		}
	}

	return company
}

func GetOwnerByID(ownerID int) *model.Owner {
	rows := databaseConnection.QueryRow("SELECT ID, FirstName, LastName, Address FROM dbo.owner WHERE ID = ?;", ownerID)

	owner := new(model.Owner)

	if err := rows.Scan(&owner.ID, &owner.FirstName, &owner.LastName, &owner.Address); err != nil {
		if err == sql.ErrNoRows {
			return nil
		} else {
			helper.HandleDatabaseError(err)
		}
	}

	return owner
}

func InitRepository() {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;",
		server, user, password, port, database)
	conn, err := sql.Open("mssql", connString)

	if err != nil {
		helper.HandleDatabaseError(err)
	}

	stmt, err := conn.Prepare("select 1, 'abc'")

	if err != nil {
		helper.HandleDatabaseError(err)
	}
	defer stmt.Close()

	databaseConnection = conn
}
