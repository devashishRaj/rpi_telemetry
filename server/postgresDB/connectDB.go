package postgresDB

import (
	"database/sql"
	"fmt"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "rospi"
	password = "rpitele"
	dbname   = "rospi"
)

var db *sql.DB

func ConnectDB() *sql.DB {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {

		fmt.Println("error in ConnectDB")
		panic(err)
	} else {
		return db
	}
}
func CloseDB() {
	if db != nil {
		db.Close()
		fmt.Println("DataBase closed")
	}
}
