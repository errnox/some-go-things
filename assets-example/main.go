package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"

	// rice "github.com/GeertJohan/go.rice"
	_ "github.com/mattn/go-sqlite3"
)

type Person struct {
	PersonId int    `json:"id" name:"ID"`
	Name     string `json:"name" name:"Name"`
	Age      int    `json:"age" name:"Age"`
}

var DB *sql.DB
var err error

func init() {
	fmt.Println("Run init")

	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	home := user.HomeDir
	dst := home + "/assetstest.sqlite"
	CopyFile("./db/db.sqlite", dst)

	// DB, err = sql.Open("sqlite3", "./db/db.sqlite")
	DB, err = sql.Open("sqlite3", dst)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("Run main")

	// As usual
	//
	// log.Fatal(http.ListenAndServe(":4321",
	// 	http.FileServer(http.Dir("./public"))))

	// Use `go.rice'.
	//
	// http.Handle("/", http.FileServer(
	// 	rice.MustFindBox("public").HTTPBox()))
	// http.ListenAndServe(":4321", nil)

	rows, err := DB.Query("SELECT * FROM Person;")
	people := []*Person{}
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		person := &Person{}
		rows.Scan(&person.PersonId,
			&person.Name,
			&person.Age)
		people = append(people, person)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	rows.Close()
	for _, person := range people {
		fmt.Println(person.PersonId, person.Name, person.Age)
	}
	rows.Close()

	defer DB.Close()
}

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(out, in)
	clerr := out.Close()
	return clerr
}

// func CopyDB(data []byte, dst string) err {
// 	out, err := os.Create(dst)
// 	if err != nil {
// 		return err
// 	}
// 	defer out.Close()
// 	io.Copy(data, out)
// 	clerr := out.Close()
// 	return clerr
// }
