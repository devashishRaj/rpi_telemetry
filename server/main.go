package main

import (
	//"bytes"

	"database/sql"
	"encoding/json"
	"fmt"

	//"io"
	"net/http"
	//"time"
	"log"

	"github.com/lib/pq"
	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
)

// database credentials
const (
	host     = "localhost"
	port     = 5432
	user     = "rospi"
	password = "rpitele"
	dbname   = "rospi"
)

// struct to store json data
type SystemInfo struct {
	HardwareID  string  `json:"HardwareID"`
	CPUuserLoad float64 `json:"CPUuserLoad"`
	CPUidle     float64 `json:"CPUidle"`
	TotalMemory int64   `json:"TotalMemory"`
	FreeMemory  int64   `json:"FreeMemory"`
	IP          string  `json:"IP"`
	Temperature string  `json:"Temperature"`
	TimeStamp   string  `json:"TimeStamp"`
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func ReceiveJSON(w http.ResponseWriter, r *http.Request) {
	var jsonData SystemInfo
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Received Info: %+v\n", jsonData)
	w.WriteHeader(http.StatusOK)

	DBerr := updateDB(jsonData)
	if DBerr != nil {
		return
	}
}

func updateDB(jsonData SystemInfo) error {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	defer db.Close()

	_, err = db.Exec(`
	INSERT INTO telemetry.rpib (HardwareID, CPUuserLoad, CPUidle, TotalMemory, FreeMemory, IP, 
								Temperature, TimeStamp)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		jsonData.HardwareID, jsonData.CPUuserLoad, jsonData.CPUidle,
		jsonData.TotalMemory, jsonData.FreeMemory, jsonData.IP,
		jsonData.Temperature, pq.QuoteLiteral(jsonData.TimeStamp))

	CheckError(err)
	fmt.Println("Data inserted successfully!")
	return nil
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/rpi", ReceiveJSON).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
