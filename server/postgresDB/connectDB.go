package postgresDB

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/spf13/viper"
)

var db *sql.DB

func ReadConfig() {
	viper.AddConfigPath("./local/.configs")
	viper.SetConfigName("config") // Register config file name (no extension)
	viper.SetConfigType("json")   // Look for specific type
	viper.ReadInConfig()
	viper.WatchConfig()
}

func ConnectDB() *sql.DB {
	ReadConfig()
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		viper.Get("postgresDB.host"), viper.Get("postgresDB.port"), viper.Get("postgresDB.user"),
		viper.Get("postgresDB.password"), viper.Get("postgresDB.dbname"),
		viper.Get("postgresDB.sslmode"))

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {

		fmt.Println("error in ConnectDB")
		log.Fatalln(err)
	} else {
		return db
	}
	return nil
}
func CloseDB() {
	if db != nil {
		db.Close()
	}
}
