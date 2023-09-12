package main

import (
	"server/jsonHandler"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	//
	jsonHandler.ReceiveJSON()

}
