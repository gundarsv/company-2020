package repository

import (
	"company-api/src/helper"
	"company-api/src/model"
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
	"os"
	"strings"
)

var (
	server             = os.Getenv("DATABASE_IP")
	user               = os.Getenv("DATABASE_USER")
	password           = os.Getenv("DATABASE_PASSWORD")
	database           = os.Getenv("DATABASE_DATABASE")
	databaseConnection *sql.DB
)

func GetAllOwners() (*helper.DatabaseResponse, []*model.Owner) {
	rows, err := databaseConnection.Query("SELECT ID, FirstName, LastName, Address FROM owner;")

	if err != nil {
		return helper.NewDatabaseResponse(true, "Something went wrong", err.Error(), helper.RiskDatabaseError), nil
	}
	defer rows.Close()

	var owners []*model.Owner

	for rows.Next() {
		o := new(model.Owner)
		err := rows.Scan(&o.ID, &o.FirstName, &o.LastName, &o.Address)

		if err != nil {
			return helper.NewDatabaseResponse(true, "Something went wrong", err.Error(), helper.RiskDatabaseError), nil
		}

		owners = append(owners, o)
	}

	return helper.NewDatabaseResponse(false, "Owners returned", "Owner returned successfully", helper.NoDatabaseError), owners
}

func GetAllCompanies() (*helper.DatabaseResponse, []*model.Company) {
	rows, err := databaseConnection.Query("SELECT ID, Name, Address, City, Country, COALESCE(Email, ''), COALESCE(PhoneNumber, '') FROM company;")

	if err != nil {
		return helper.NewDatabaseResponse(true, "Something went wrong", err.Error(), helper.RiskDatabaseError), nil
	}
	defer rows.Close()

	var companies []*model.Company

	for rows.Next() {
		c := new(model.Company)
		err := rows.Scan(&c.ID, &c.Name, &c.Address, &c.City, &c.Country, &c.PhoneNumber, &c.Email)

		if err != nil {
			return helper.NewDatabaseResponse(true, "No companies found", err.Error(), helper.NotFoundDatabaseError), nil
		}

		if response := addOwnersToCompany(c); response.IsError {
			return helper.NewDatabaseResponse(response.IsError, response.Message, response.LogMessage, response.MessageCode), nil
		}

		companies = append(companies, c)
	}

	return helper.NewDatabaseResponse(false, "Companies returned", "Companies returned successfully", helper.NoDatabaseError), companies
}

func addOwnersToCompany(company *model.Company) *helper.DatabaseResponse {
	rows, err := databaseConnection.Query("SELECT * \n"+
		"FROM owner \n"+
		"WHERE id IN (SELECT ownerID FROM company_owner WHERE companyID = ?);", company.ID)

	if err != nil {
		return helper.NewDatabaseResponse(true, "Something went wrong", err.Error(), helper.RiskDatabaseError)
	}

	for rows.Next() {
		owner := new(model.Owner)

		err := rows.Scan(&owner.ID, &owner.FirstName, &owner.LastName, &owner.Address)

		if err != nil {
			return helper.NewDatabaseResponse(true, "Something went wrong", err.Error(), helper.RiskDatabaseError)
		}

		company.AddOwner(*owner)
	}

	return helper.NewDatabaseResponse(false, "Owner added", "Owner added to company", helper.NoDatabaseError)
}

func UpdateCompany(id int, company model.Company) (*helper.DatabaseResponse, *model.Company) {
	queryString := "UPDATE [company] set "
	var queryArgs []interface{}

	if company.Name != helper.NilString {
		queryString = queryString + "Name = coalesce(?, Name),"
		queryArgs = append(queryArgs, company.Name)
	}

	if company.Address != helper.NilString {
		queryString = queryString + "Address = coalesce(?, Address),"
		queryArgs = append(queryArgs, company.Address)
	}

	if company.City != helper.NilString {
		queryString = queryString + "City = coalesce(?, City),"
		queryArgs = append(queryArgs, company.City)
	}

	if company.Country != helper.NilString {
		queryString = queryString + "Country = coalesce(?, Country),"
		queryArgs = append(queryArgs, company.Country)
	}

	if company.Email != helper.NilString {
		queryString = queryString + "Email = coalesce(?, Email),"
		queryArgs = append(queryArgs, company.Email)
	}

	if company.PhoneNumber != helper.NilString {
		queryString = queryString + "PhoneNumber = coalesce(?, PhoneNumber),"
		queryArgs = append(queryArgs, company.PhoneNumber)
	}

	if len(queryArgs) <= 0 {
		return helper.NewDatabaseResponse(true, "Please add parameters", "No parameters were added correctly", helper.NoParametersDatabaseError), nil
	}

	queryString = helper.TrimSuffix(queryString, ",")

	queryString = queryString + " WHERE ID = ?"
	queryArgs = append(queryArgs, id)

	_, err := databaseConnection.Exec(queryString, queryArgs...)

	if err != nil {
		return helper.NewDatabaseResponse(true, "Something went wrong", err.Error(), helper.RiskDatabaseError), nil
	}

	return GetCompanyByID(id)
}

func UpdateOwner(id int, owner model.Owner) (*helper.DatabaseResponse, *model.Owner) {
	queryString := "UPDATE [owner] set "
	var queryArgs []interface{}

	if owner.FirstName != helper.NilString {
		queryString = queryString + "FirstName = coalesce(?, FirstName),"
		queryArgs = append(queryArgs, owner.FirstName)
	}

	if owner.LastName != helper.NilString {
		queryString = queryString + "LastName = coalesce(?, LastName),"
		queryArgs = append(queryArgs, owner.LastName)
	}

	if owner.Address != helper.NilString {
		queryString = queryString + "Address = coalesce(?, Address),"
		queryArgs = append(queryArgs, owner.Address)
	}

	if len(queryArgs) <= 0 {
		return helper.NewDatabaseResponse(true, "Please add parameters", "No parameters were added correctly", helper.NoParametersDatabaseError), nil
	}

	queryString = helper.TrimSuffix(queryString, ",")

	queryString = queryString + " WHERE ID = ?"
	queryArgs = append(queryArgs, id)

	_, err := databaseConnection.Exec(queryString, queryArgs...)

	if err != nil {
		return helper.NewDatabaseResponse(true, "Something went wrong", err.Error(), helper.RiskDatabaseError), nil
	}

	return GetOwnerByID(id)
}

func DeleteCompany(id int) *helper.DatabaseResponse {
	rows, err := databaseConnection.Exec("DELETE FROM [company] WHERE ID = ?;", id)

	if err != nil {
		return helper.NewDatabaseResponse(true, "Something went wrong", err.Error(), helper.RiskDatabaseError)
	}

	number, err := rows.RowsAffected()

	if err != nil {
		return helper.NewDatabaseResponse(true, "Something went wrong", err.Error(), helper.RiskDatabaseError)
	}

	if number <= 0 {
		helper.NewDatabaseResponse(true, "No company with id:"+string(id), "No company with id:"+string(id)+" was found", helper.NotFoundDatabaseError)
	}

	return helper.NewDatabaseResponse(false, "Company with id:"+string(id)+" deleted", "Company with id:"+string(id)+" deleted", helper.NoDatabaseError)
}

func DeleteOwner(id int) *helper.DatabaseResponse {
	rows, err := databaseConnection.Exec("DELETE FROM [owner] WHERE ID = ?;", id)

	if err != nil {
		return helper.NewDatabaseResponse(true, "Something went wrong", err.Error(), helper.RiskDatabaseError)
	}

	number, err := rows.RowsAffected()

	if err != nil {
		return helper.NewDatabaseResponse(true, "Something went wrong", err.Error(), helper.RiskDatabaseError)
	}

	if number <= 0 {
		helper.NewDatabaseResponse(true, "No owner with id:"+string(id), "No owner with id:"+string(id)+"was found", helper.NotFoundDatabaseError)
	}

	return helper.NewDatabaseResponse(false, "Owner with id:"+string(id)+" deleted", "Owner with id:"+string(id)+" was deleted", helper.NoDatabaseError)
}

func DeleteOwnerFromCompany(companyID int, ownerID int) *helper.DatabaseResponse {
	rows, err := databaseConnection.Exec("DELETE FROM [company_owner] WHERE companyID = ? and ownerID = ?;", companyID, ownerID)

	if err != nil {
		return helper.NewDatabaseResponse(true, "Something went wrong", err.Error(), helper.RiskDatabaseError)
	}

	number, err := rows.RowsAffected()

	if err != nil {
		return helper.NewDatabaseResponse(true, "Something went wrong", err.Error(), helper.RiskDatabaseError)
	}

	if number <= 0 {
		helper.NewDatabaseResponse(true, "No company with id:"+string(companyID)+"or no owner with id:"+string(ownerID), "No company with id:"+string(companyID)+"or no owner with id:"+string(ownerID), helper.NoDatabaseError)
	}

	return helper.NewDatabaseResponse(false, "Owner with id:"+string(ownerID)+"was deleted from company with id:"+string(companyID), "Owner with id:"+string(ownerID)+"was deleted from company with id:"+string(companyID), helper.NoDatabaseError)
}

func CreateCompany(company model.Company) (*helper.DatabaseResponse, *model.Company) {
	rows := databaseConnection.QueryRow("INSERT INTO [company]([Name],[Address],[Country],[City],[Email],[PhoneNumber]) VALUES (?, ?, ?, ?, ?, ?); select ID = convert(bigint, SCOPE_IDENTITY());",
		&company.Name, &company.Address, &company.Country, &company.City, &company.Email, &company.PhoneNumber)

	id := new(int)

	err := rows.Scan(&id)

	if err != nil {
		return helper.NewDatabaseResponse(true, "Something went wrong", err.Error(), helper.RiskDatabaseError), nil
	}

	return GetCompanyByID(*id)
}

func CreateOwner(owner model.Owner) (*helper.DatabaseResponse, *model.Owner) {
	rows := databaseConnection.QueryRow("INSERT INTO [owner]([FirstName],[LastName],[Address]) VALUES (?, ?, ?); select ID = convert(bigint, SCOPE_IDENTITY());",
		&owner.FirstName, &owner.LastName, &owner.Address)

	id := new(int)

	err := rows.Scan(&id)

	if err != nil {
		return helper.NewDatabaseResponse(true, "Something went wrong", err.Error(), helper.RiskDatabaseError), nil
	}

	return GetOwnerByID(*id)
}

func AddOwnerToCompany(companyID int, ownerID int) (*helper.DatabaseResponse, *model.Company) {
	_, err := databaseConnection.Exec("INSERT INTO [company_owner]([CompanyID],[OwnerID]) VALUES (?, ?);", companyID, ownerID)

	if err != nil {
		if strings.Contains(err.Error(), "UQ_CompanyID_OwnerID") {
			return helper.NewDatabaseResponse(true, "Company already has this owner", err.Error(), helper.AlreadyExistsDatabaseError), nil
		}

		if strings.Contains(err.Error(), "FK_company_owner_company_owner") {
			return helper.NewDatabaseResponse(true, "Company does not exist", err.Error(), helper.NotFoundDatabaseError), nil
		}

		if strings.Contains(err.Error(), "FK_company_owner_owner") {
			return helper.NewDatabaseResponse(true, "Owner does not exist", err.Error(), helper.NotFoundDatabaseError), nil
		}

		return helper.NewDatabaseResponse(true, "Something went wrong", err.Error(), helper.RiskDatabaseError), nil
	}

	return GetCompanyByID(companyID)
}

func GetCompanyByID(companyID int) (*helper.DatabaseResponse, *model.Company) {
	rowsCompany := databaseConnection.QueryRow("SELECT ID, Name, Address, City, Country, COALESCE(Email, ''), COALESCE(PhoneNumber, '') FROM company WHERE ID = ?;", companyID)

	company := new(model.Company)

	err := rowsCompany.Scan(&company.ID, &company.Name, &company.Address, &company.City, &company.Country, &company.Email, &company.PhoneNumber)

	if err == sql.ErrNoRows {
		return helper.NewDatabaseResponse(true, "Company does not exist", sql.ErrNoRows.Error(), helper.NotFoundDatabaseError), nil
	}

	if err != nil {
		return helper.NewDatabaseResponse(true, "Something went wrong", err.Error(), helper.RiskDatabaseError), nil
	}

	addOwnersToCompany(company)

	return helper.NewDatabaseResponse(false, "Company returned successfully", "Company returned successfully", helper.RiskDatabaseError), company
}

func GetOwnerByID(ownerID int) (*helper.DatabaseResponse, *model.Owner) {
	rows := databaseConnection.QueryRow("SELECT ID, FirstName, LastName, Address FROM owner WHERE ID = ?;", ownerID)

	owner := new(model.Owner)

	err := rows.Scan(&owner.ID, &owner.FirstName, &owner.LastName, &owner.Address)

	if err == sql.ErrNoRows {
		return helper.NewDatabaseResponse(true, "Owner does not exist", sql.ErrNoRows.Error(), helper.NotFoundDatabaseError), nil
	}

	if err != nil {
		return helper.NewDatabaseResponse(true, "Something went wrong", err.Error(), helper.RiskDatabaseError), nil
	}

	return helper.NewDatabaseResponse(false, "Owner returned successfully", "Owner returned successfully", helper.NoDatabaseError), owner
}

//func createDatabase() {
//	databaseCreationQuery := "CREATE DATABASE companyDB;"
//
//	databaseCreationResult, err := databaseConnection.Exec(databaseCreationQuery)
//
//	if err != nil {
//		helper.HandleDatabaseError(err)
//	}
//
//	databaseCreationRowsAffected, err := databaseCreationResult.RowsAffected()
//
//	if err != nil {
//		helper.HandleDatabaseError(err)
//	}
//
//	log.Println("Database created " + string(databaseCreationRowsAffected))
//
//	tableCreationQuery :=
//		"CREATE TABLE [companyDB].[dbo].[company]( \n" +
//			"[ID] [int] IDENTITY(1,1) NOT NULL, \n" +
//			"[Name] [varchar](50) NOT NULL, \n" +
//			"[Address] [varchar](50) NOT NULL, \n" +
//			"[Country] [varchar](50) NOT NULL, \n" +
//			"[City] [varchar](50) NOT NULL, \n" +
//			"[Email] [varchar](50) NULL, \n" +
//			"[PhoneNumber] [varchar](50) NULL, \n" +
//			"CONSTRAINT [PK_company] PRIMARY KEY CLUSTERED ([ID] ASC \n" +
//			")WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY] \n" +
//			") ON [PRIMARY];" +
//			"CREATE TABLE [companyDB].[dbo].[company_owner]( \n" +
//			"[ID] [int] IDENTITY(1,1) NOT NULL, \n" +
//			"[CompanyID] [int] NOT NULL, \n" +
//			"[OwnerID] [int] NOT NULL, \n" +
//			"CONSTRAINT [PK_company_owner] PRIMARY KEY CLUSTERED ([ID] ASC \n" +
//			")WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY] \n" +
//			") ON [PRIMARY]; \n" +
//			"CREATE TABLE [companyDB].[dbo].[owner]( \n" +
//			"[ID] [int] IDENTITY(1,1) NOT NULL, \n" +
//			"[FirstName] [varchar](50) NOT NULL, \n" +
//			"[LastName] [varchar](50) NOT NULL, \n" +
//			"[Address] [varchar](50) NOT NULL, \n" +
//			"CONSTRAINT [PK_owner] PRIMARY KEY CLUSTERED ([ID] ASC \n" +
//			")WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY] \n" +
//			") ON [PRIMARY]; \n"
//
//	tableCreationResult, err := databaseConnection.Exec(tableCreationQuery)
//
//	if err != nil {
//		helper.HandleDatabaseError(err)
//	}
//
//	tableCreationRowsAffected, err := tableCreationResult.RowsAffected()
//
//	if err != nil {
//		helper.HandleDatabaseError(err)
//	}
//
//	log.Println("Tables created " + string(tableCreationRowsAffected))
//
//	keyRelationshipCreationQuery :=
//		"ALTER TABLE [companyDB].[dbo].[company_owner] ADD  CONSTRAINT [UQ_CompanyID_OwnerID] UNIQUE NONCLUSTERED \n" +
//			"([CompanyID] ASC, \n" +
//			"[OwnerID] ASC \n" +
//			")WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, SORT_IN_TEMPDB = OFF, IGNORE_DUP_KEY = OFF, ONLINE = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]; \n" +
//			"ALTER TABLE [companyDB].[dbo].[company_owner]  WITH CHECK ADD  CONSTRAINT [FK_company_owner_company_owner] FOREIGN KEY([CompanyID]) \n" +
//			"REFERENCES [companyDB].[dbo].[company] ([ID]) ON DELETE CASCADE; \n" +
//			"ALTER TABLE [companyDB].[dbo].[company_owner] CHECK CONSTRAINT [FK_company_owner_company_owner]; \n" +
//			"ALTER TABLE [companyDB].[dbo].[company_owner]  WITH CHECK ADD  CONSTRAINT [FK_company_owner_owner] FOREIGN KEY([OwnerID]) \n" +
//			"REFERENCES [companyDB].[dbo].[owner] ([ID]) ON DELETE CASCADE; \n" +
//			"ALTER TABLE [companyDB].[dbo].[company_owner] CHECK CONSTRAINT [FK_company_owner_owner]"
//
//	keyRelationshipCreationResult, err := databaseConnection.Exec(keyRelationshipCreationQuery)
//
//	if err != nil {
//		helper.HandleDatabaseError(err)
//	}
//
//	keyRelationshipCreationRowsAffected, err := keyRelationshipCreationResult.RowsAffected()
//
//	if err != nil {
//		helper.HandleDatabaseError(err)
//	}
//
//	log.Println("Key relationships added " + string(keyRelationshipCreationRowsAffected))
//}

//func checkIfDatabaseExists() {
//	row := databaseConnection.QueryRow("SELECT COUNT(name) FROM master.dbo.sysdatabases WHERE name = ?;", "companyDB")
//
//	exists := false
//	if err := row.Scan(&exists); err != nil {
//		helper.HandleDatabaseError(err)
//	}
//
//	if !exists {
//		log.Println("Database doesn't exist")
//		createDatabase()
//	}
//}

func InitRepository() {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s",
		server, user, password, database)
	conn, err := sql.Open("mssql", connString)

	if err != nil {
		log.Fatalln(err)
	}

	if err := conn.Ping(); err != nil {
		log.Fatalln(err)
	}

	databaseConnection = conn
}
