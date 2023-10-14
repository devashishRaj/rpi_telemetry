package main

import (
	"devashishRaj/rpi_telemetry/server/jsonHandler"
	"devashishRaj/rpi_telemetry/server/postgresDB"
	"log"
)

func main() {
	log.Println("starting server")
	postgresDB.ConnectDB()
	jsonHandler.ReceiveJSON()
	defer postgresDB.G_dbpool.Close()
}
