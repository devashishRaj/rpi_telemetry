package postgresDB

import (
	"context"
	"fmt"

	handle "devashishRaj/rpi_telemetry/server/handleError"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/spf13/viper"
)

//var G_dbpool *pgxpool.Pool

func ReadConfig() {
	viper.AddConfigPath("./local/.configs")
	viper.SetConfigName("config") // Register config file name (no extension)
	viper.SetConfigType("json")   // Look for specific type
	err := viper.ReadInConfig()
	handle.CheckError("unable to read config , func: ReadConfig", err)
}

func ConnectDB() *pgxpool.Pool {
	ReadConfig()
	// connection string
	psqlConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		viper.Get("postgresDB.host"), viper.Get("postgresDB.port"), viper.Get("postgresDB.user"),
		viper.Get("postgresDB.password"), viper.Get("postgresDB.dbname"),
		viper.Get("postgresDB.sslmode"))

	// open database
	db, err := pgxpool.New(context.Background(), psqlConnStr)
	handle.CheckError("Unable to connect to DB , func: ConnectDB", err)
	handle.CheckError("Database connection health", err)
	return db
}
