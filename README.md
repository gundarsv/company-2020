# Company 2020
Company 2020 is an application made to organize companies and company owners.

Requirements:
* [Docker](https://www.docker.com/get-started)
* [Microsoft SQL Server*]()

To run the application:
1. Clone the repository
2. Create a *.env* file in **./company-2020**  for all necessary environment variables**.
3. Run ```CREATE_DATABASE.sql``` from your MS SQL Server.
4. In **./company-2020** run ```docker-compose up```.
5. Application is now running locally to the **port** you specified in *.env*.

*\*Can be 'dockerized' MS SQL Server, learn how to do that [here](https://docs.microsoft.com/en-us/sql/linux/quickstart-install-connect-docker?view=sql-server-ver15&pivots=cs1-powershell).*  
*\*\*.env file should look like this ```DATABASE={DatabaseName}
                                        HOST={DatabaseServerURI} 
                                        USER={DatabasUser}  
                                        PASSWORD={DatabasePassword}
                                        PORT={PortForApplication}```*


---
### Company WEB

Frontend of the Company 2020 Application where you can manage companies and company owners.  

Tools used: 
* React 16.12
* Typescript 3.7.5
* Webpack
* Material-UI

---
### Company API
Backend of the application which connects to the MS SQL Server.
Company API is written in Go programming language. [Read about Go here](). To connect to MS SQL Server it uses [denisenkom/go-mssqldb](https://github.com/denisenkom/go-mssqldb) driver. For HTTP Request routing [gorilla/mux](https://github.com/gorilla/mux)
 is used. 

---
### HTTP requests to Company API
You need to set the **{url}** either to **localhost** if you are running locally or **company2020.herokuapp.com** to run against deployed application.

#### GET Requests
* Do  ```curl --header "Content-Type: application/json" --request GET {url}/api/company``` to return a list of all companies in database.  
  
* Do  ```curl --header "Content-Type: application/json" --request GET {url}/api/owner``` to return a list of all owners in database.  
  
 * Do  ```curl --header "Content-Type: application/json" --request GET {url}/api/company/{id}``` where *id* is the ID of requested company, to get the requested company.
          
 * Do  ```curl --header "Content-Type: application/json" --request GET {url}/api/owner/{id}``` where *id* is the ID of requested owner, to get the requested owner.  
 
 #### POST Requests
 **Values *{required}* are the values required to create the entity.**  
**Values *{optional}* are the values that are not required to create the entity.**
 
* Do  ```curl --header "Content-Type: application/json" 
             --request POST
             --data '{"Name":"{required}","Address":"{required}","City":"{required}","Country":"{required}","Email":"{optional}","PhoneNumber":"{optional}"}'
             {url}/api/company``` to create a new company.
             
* Do  ```curl --header "Content-Type: application/json" 
             --request POST
             --data '{"FirstName":"{required}","LastName":"{required}","Address":"{required}"}'
             {url}/api/owner``` to create a new owner.               
             
 #### PUT Requests
 **Values *{optional}* are the values that are not required to update the entity.**
* Do  ```curl --header "Content-Type: application/json" 
              --request PUT
              --data '{"Name":"{optional}","Address":"{optional}","City":"{optional}","Country":"{optional}","Email":"{optional}","PhoneNumber":"{optional}"}'
              {url}/api/company/{id}``` where *id* is the ID of company you want to update, to update the company.
              
* Do  ```curl --header "Content-Type: application/json" 
              --request PUT
              --data '{"FirstName":"{optional}","LastName":"{optional}","Address":"{optional}"}'
              {url}/api/owner/{id}``` where *id* is the ID of owner you want to update, to update the owner.
              
* Do  ```curl --header "Content-Type: application/json" 
                            --request PUT
                            {url}/api/company/{companyId}/owner/{ownerId}``` where *companyId* is the ID of company you want to add owner to and where *ownerId* is the ID of owner you want to add to the company, to add owner to company.
  
 #### DELETE Requests
 
* Do  ```curl --header "Content-Type: application/json" 
               --request DELETE
               {url}/api/owner/{id}``` where *id* is the ID of owner you want to delete, to delete the owner.
               
* Do  ```curl --header "Content-Type: application/json" 
                --request DELETE
                {url}/api/company/{id}``` where *id* is the ID of company you want to delete, to delete the company.
                
* Do  ```curl --header "Content-Type: application/json" 
                            --request DELETE
                            {url}/api/company/{companyId}/owner/{ownerId}``` where *companyId* is the ID of company you want to delete the owner from to and where *ownerId* is the ID of owner you want to delete from the company, to delete owner from company.