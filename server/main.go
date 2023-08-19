package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
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

func main() {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	defer db.Close()

	// loop to receive updated value using HTTP GET
	for range time.Tick(time.Second * 9) {
		raspberryPiIP := "192.168.1.4:8080"

		response, err := http.Get("http://" + raspberryPiIP + "/send-json")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer response.Body.Close()
		//reads output from GET
		var bodyBuffer bytes.Buffer
		_, copyErr := io.Copy(&bodyBuffer, response.Body)
		CheckError(copyErr)
		body := bodyBuffer.Bytes()
		// convert json format datat to struct
		var jsonData SystemInfo
		err = json.Unmarshal(body, &jsonData)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Printf("System: %+v\n", jsonData)
		// insert values into db
		_, err = db.Exec(`
		INSERT INTO telemetry.rpib (HardwareID, CPUuserLoad, CPUidle, TotalMemory, FreeMemory, IP, Temperature, TimeStamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			jsonData.HardwareID, jsonData.CPUuserLoad, jsonData.CPUidle,
			jsonData.TotalMemory, jsonData.FreeMemory, jsonData.IP,
			jsonData.Temperature, pq.QuoteLiteral(jsonData.TimeStamp))

		CheckError(err)
		fmt.Println("Data inserted successfully!")

	}
}
