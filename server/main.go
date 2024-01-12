package main

import (
	"devashishRaj/rpi_telemetry/server/jsonHandler"
	"log"
)

func main() {
	log.Println("starting server")
	jsonHandler.ReceiveJSON()

}
