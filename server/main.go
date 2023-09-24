package main

import (
	"log"
	"server/jsonHandler"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	//
	log.Println("starting server")
	jsonHandler.ReceiveJSON()

}
