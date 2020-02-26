package repository

import (
	"company-api/src/helper"
	"company-api/src/model"
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
	"os"
)

var (
	server             = os.Getenv("DATABASE_IP")
	port               = os.Getenv("DATABASE_PORT")
	user               = "sa"
	password           = "Secret!Secret"
	databaseConnection *sql.DB
)

func GetAllOwners() []*model.Owner {
	rows, err := databaseConnection.Query("SELECT ID, FirstName, LastName, Address FROM companyDB.dbo.owner;")

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
	rows, err := databaseConnection.Query("SELECT ID, Name, Address, City, Country, COALESCE(Email, ''), COALESCE(PhoneNumber, '') FROM companyDB.dbo.company;")

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
	rows := databaseConnection.QueryRow("SELECT ID, Name, Address, City, Country, COALESCE(Email, ''), COALESCE(PhoneNumber, '') FROM companyDB.dbo.company WHERE ID = ?;", companyID)

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
	rows := databaseConnection.QueryRow("SELECT ID, FirstName, LastName, Address FROM companyDB.dbo.owner WHERE ID = ?;", ownerID)

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

func createDatabase() {
	databaseCreationQuery := "CREATE DATABASE companyDB;"

	databaseCreationResult, err := databaseConnection.Exec(databaseCreationQuery)

	if err != nil {
		helper.HandleDatabaseError(err)
	}

	databaseCreationRowsAffected, err := databaseCreationResult.RowsAffected()

	if err != nil {
		helper.HandleDatabaseError(err)
	}

	log.Println("Database created " + string(databaseCreationRowsAffected))

	tableCreationQuery :=
		"CREATE TABLE [companyDB].[dbo].[company]( \n" +
			"[ID] [int] IDENTITY(1,1) NOT NULL, \n" +
			"[Name] [varchar](50) NOT NULL, \n" +
			"[Address] [varchar](50) NOT NULL, \n" +
			"[Country] [varchar](50) NOT NULL, \n" +
			"[City] [varchar](50) NOT NULL, \n" +
			"[Email] [varchar](50) NULL, \n" +
			"[PhoneNumber] [varchar](50) NULL, \n" +
			"CONSTRAINT [PK_company] PRIMARY KEY CLUSTERED ([ID] ASC \n" +
			")WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY] \n" +
			") ON [PRIMARY];" +
			"CREATE TABLE [companyDB].[dbo].[company_owner]( \n" +
			"[ID] [int] IDENTITY(1,1) NOT NULL, \n" +
			"[CompanyID] [int] NOT NULL, \n" +
			"[OwnerID] [int] NOT NULL, \n" +
			"CONSTRAINT [PK_company_owner] PRIMARY KEY CLUSTERED ([ID] ASC \n" +
			")WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY] \n" +
			") ON [PRIMARY]; \n" +
			"CREATE TABLE [companyDB].[dbo].[owner]( \n" +
			"[ID] [int] IDENTITY(1,1) NOT NULL, \n" +
			"[FirstName] [varchar](50) NOT NULL, \n" +
			"[LastName] [varchar](50) NOT NULL, \n" +
			"[Address] [varchar](50) NOT NULL, \n" +
			"CONSTRAINT [PK_owner] PRIMARY KEY CLUSTERED ([ID] ASC \n" +
			")WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY] \n" +
			") ON [PRIMARY]; \n"

	tableCreationResult, err := databaseConnection.Exec(tableCreationQuery)

	if err != nil {
		helper.HandleDatabaseError(err)
	}

	tableCreationRowsAffected, err := tableCreationResult.RowsAffected()

	if err != nil {
		helper.HandleDatabaseError(err)
	}

	log.Println("Tables created " + string(tableCreationRowsAffected))

	keyRelationshipCreationQuery :=
		"ALTER TABLE [companyDB].[dbo].[company_owner] ADD  CONSTRAINT [UQ_CompanyID_OwnerID] UNIQUE NONCLUSTERED \n" +
			"([CompanyID] ASC, \n" +
			"[OwnerID] ASC \n" +
			")WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, SORT_IN_TEMPDB = OFF, IGNORE_DUP_KEY = OFF, ONLINE = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]; \n" +
			"ALTER TABLE [companyDB].[dbo].[company_owner]  WITH CHECK ADD  CONSTRAINT [FK_company_owner_company_owner] FOREIGN KEY([CompanyID]) \n" +
			"REFERENCES [companyDB].[dbo].[company] ([ID]); \n" +
			"ALTER TABLE [companyDB].[dbo].[company_owner] CHECK CONSTRAINT [FK_company_owner_company_owner]; \n" +
			"ALTER TABLE [companyDB].[dbo].[company_owner]  WITH CHECK ADD  CONSTRAINT [FK_company_owner_owner] FOREIGN KEY([OwnerID]) \n" +
			"REFERENCES [companyDB].[dbo].[owner] ([ID]); \n" +
			"ALTER TABLE [companyDB].[dbo].[company_owner] CHECK CONSTRAINT [FK_company_owner_owner]"

	keyRelationshipCreationResult, err := databaseConnection.Exec(keyRelationshipCreationQuery)

	if err != nil {
		helper.HandleDatabaseError(err)
	}

	keyRelationshipCreationRowsAffected, err := keyRelationshipCreationResult.RowsAffected()

	if err != nil {
		helper.HandleDatabaseError(err)
	}

	log.Println("Key relationships added " + string(keyRelationshipCreationRowsAffected))
}

func checkIfDatabaseExists() {
	row := databaseConnection.QueryRow("SELECT COUNT(name) FROM master.dbo.sysdatabases WHERE name = 'companyDB'")

	exists := false
	if err := row.Scan(&exists); err != nil {
		helper.HandleDatabaseError(err)
	}

	if !exists {
		log.Println("Database doesn't exist")
		createDatabase()
	}
}

func InitRepository() {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;",
		server, user, password, port)
	conn, err := sql.Open("mssql", connString)

	if err != nil {
		helper.HandleDatabaseError(err)
	}

	if err := conn.Ping(); err != nil {
		helper.HandleDatabaseError(err)
	}

	databaseConnection = conn

	checkIfDatabaseExists()
}
