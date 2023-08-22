package main

import (
	connect "server/postgresDB"

	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/rpi", connect.ReceiveJSON).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))

}
