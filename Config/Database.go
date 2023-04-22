package Config

import (
	"database/sql"
	"fmt"
)

// ConnectToDB connects to the database --> orders
func ConnectToDB() *sql.DB {
	db, err := sql.Open("mysql", "root:Ranjith@tcp(localhost:3306)/sachindb")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println("Connected to DB Successfully....... ")
	return db
}

// NewTable creates new table if the table not exist
func NewTable() {
	db := ConnectToDB()
	defer db.Close()
	_, err := db.Query("CREATE TABLE IF NOT EXISTS users(Name varchar(20) UNIQUE NOT NULL, Username varchar(20) NOT NULL, Email varchar(20) NOT NULL, Password varchar(20) NOT NULL)")
	if err != nil {
		fmt.Println(err)
	}
	_, e := db.Query("CREATE TABLE IF NOT EXISTS todo(ID int(5) UNIQUE NOT NULL, Title varchar(20) NOT NULL, Description varchar(20))")
	if e != nil {
		fmt.Println(e)
	}
}
