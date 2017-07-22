package datastore

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"

	_ "github.com/mattn/go-sqlite3"
)

type Customer struct {
	CustomerId   int    `json:"id" name:"ID"`
	FirstName    string `json:"firstName" name:"First Name"`
	LastName     string `json:"lastName" name:"Last Name"`
	Company      string `json:"company" name:"Company"`
	Address      string `json:"address" name:"Address"`
	City         string `json:"city" name:"City"`
	State        string `json:"state" name:"State"`
	Country      string `json:"country" name:"Country"`
	PostalCode   string `json:"postalCode" name:"Postal Code"`
	Phone        string `json:"phone" name:"Phone"`
	Fax          string `json:"fax" name:"Fax"`
	Email        string `json:"email" name:"Email"`
	NumData      int    `json:"numData" name:"Num Data"`
	SupportRepId int    `json:"supportRepID" name:"Support Rep ID"`
}

var DB *sql.DB
var err error

func init() {
	DB, err = sql.Open("sqlite3", "./db/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
}

func QueryCustomers() []*Customer {
	rows, err := DB.Query("SELECT * FROM Customer ORDER BY CustomerID DESC;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	customers := []*Customer{}
	x := 0.0
	for rows.Next() {
		customer := &Customer{}
		x = 0
		rows.Scan(&customer.CustomerId,
			&customer.FirstName,
			&customer.LastName,
			&customer.Company,
			&customer.Address,
			&customer.City,
			&customer.State,
			&customer.Country,
			&customer.PostalCode,
			&customer.Phone,
			&customer.Fax,
			&customer.Email,
			&customer.SupportRepId)
		v := reflect.ValueOf(*customer)
		for i := 0; i < v.NumField()-2; i++ {
			if v.Field(i).Kind() == reflect.String &&
				v.Field(i).String() == "" {
				x++
			}
		}
		customer.NumData = 100 - int(x/(float64(v.NumField()))*100)

		customers = append(customers, customer)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	rows.Close()
	return customers
}

func QueryCustomersBySearchTerm(s string) []*Customer {
	query := fmt.Sprintf("%%%s%%", s)
	rows, err := DB.Query(`SELECT *
 FROM Customer
 WHERE FirstName LIKE ?
 OR LastName LIKE ?
 OR Company LIKE ?
 OR Address LIKE ?
 OR City LIKE ?
 OR State LIKE ?
 OR Country LIKE ?
 OR Email LIKE ?
;`, query, query, query, query, query, query, query, query)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	customers := []*Customer{}
	for rows.Next() {
		customer := &Customer{}
		rows.Scan(&customer.CustomerId,
			&customer.FirstName,
			&customer.LastName,
			&customer.Company,
			&customer.Address,
			&customer.City,
			&customer.State,
			&customer.Country,
			&customer.PostalCode,
			&customer.Phone,
			&customer.Fax,
			&customer.Email,
			&customer.SupportRepId)
		customers = append(customers, customer)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	rows.Close()
	return customers
}

func QueryCustomerById(id int) *Customer {
	customer := &Customer{}

	err = DB.QueryRow("SELECT * FROM Customer WHERE CustomerId = ?",
		id).Scan(&customer.CustomerId,
		&customer.FirstName,
		&customer.LastName,
		&customer.Company,
		&customer.Address,
		&customer.City,
		&customer.State,
		&customer.Country,
		&customer.PostalCode,
		&customer.Phone,
		&customer.Fax,
		&customer.Email,
		&customer.SupportRepId)
	return customer
}

func CreateCustomer(args ...string) error {
	firstName := ""
	if len(args) >= 1 {
		firstName = args[0]
	}
	lastName := ""
	if len(args) >= 2 {
		lastName = args[1]
	}
	company := ""
	if len(args) >= 3 {
		company = args[2]
	}
	address := ""
	if len(args) >= 4 {
		address = args[3]
	}
	city := ""
	if len(args) >= 5 {
		city = args[4]
	}
	state := ""
	if len(args) >= 6 {
		state = args[5]
	}
	country := ""
	if len(args) >= 7 {
		country = args[6]
	}
	postalCode := ""
	if len(args) >= 8 {
		postalCode = args[7]
	}
	phone := ""
	if len(args) >= 9 {
		phone = args[8]
	}
	fax := ""
	if len(args) >= 10 {
		fax = args[9]
	}
	email := ""
	if len(args) >= 11 {
		email = args[10]
	}
	supportRepId := ""
	if len(args) >= 12 {
		supportRepId = args[11]
	}
	_, err := DB.Exec("INSERT INTO Customer VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		nil, firstName, lastName, company, address, city, state,
		country, postalCode, phone, fax, email, supportRepId)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not create user: %s",
			err))
	} else {
		return nil
	}
}

func DeleteCustomer(id int) error {
	_, err := DB.Exec("DELETE FROM Customer WHERE CustomerId = ?", id)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not DELETE Customer: %s",
			err))
	} else {
		return nil
	}
}

func UpdateCustomer(args ...string) error {
	id := ""
	if len(args) >= 1 {
		id = args[0]
	}
	firstName := ""
	if len(args) >= 2 {
		firstName = args[1]
	}
	lastName := ""
	if len(args) >= 3 {
		lastName = args[2]
	}
	company := ""
	if len(args) >= 4 {
		company = args[3]
	}
	address := ""
	if len(args) >= 5 {
		address = args[4]
	}
	city := ""
	if len(args) >= 6 {
		city = args[5]
	}
	state := ""
	if len(args) >= 7 {
		state = args[6]
	}
	country := ""
	if len(args) >= 8 {
		country = args[7]
	}
	postalCode := ""
	if len(args) >= 9 {
		postalCode = args[8]
	}
	phone := ""
	if len(args) >= 10 {
		phone = args[9]
	}
	fax := ""
	if len(args) >= 11 {
		fax = args[10]
	}
	email := ""
	if len(args) >= 12 {
		email = args[11]
	}
	_, err := DB.Exec(`UPDATE Customer SET
 FirstName=?,
 LastName=?,
 Company=?,
 Address=?,
 City=?,
 State=?,
 Country=?,
 PostalCode=?,
 Phone=?,
 Fax=?,
 Email=?
 WHERE CustomerId=?`, firstName, lastName, company, address, city,
		state, country, postalCode, phone, fax, email, id)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not UPDATE Customer: %s",
			err))
	} else {
		return nil
	}
}
