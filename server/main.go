package main

import (
	connect "server/jsonHandler"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	
	// router := mux.NewRouter()
	// router.HandleFunc("/rpi", connect.ReceiveJSON).Methods("POST")
	// log.Fatal(http.ListenAndServe(":8080", router))
	connect.ReceiveJSON()

}
