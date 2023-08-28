package postgresDB

import (
	"fmt"
	"log"
	dataStruct "server/dataStruct"
)

func AlertTemp(jsonData dataStruct.SystemInfo) {
	// Set custom prefix and flags to differentiate alerts
	log.SetPrefix("[ALERT] ")
	//The log.SetFlags function is used to configure the log message format to include the date,
	//time, and source file information.
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	query := `
		SELECT AVG(temperature) AS avg_temperature
		FROM telemetry.rpib
		WHERE timestamp >= NOW() - INTERVAL '60 seconds';
	`
	var avgTemperature float64
	err := db.QueryRow(query).Scan(&avgTemperature)
	if err != nil {
		fmt.Println("Query error")
	}
	CheckError(err)
	if avgTemperature > 44.5 {
		log.Printf("Average temperature within the last 30 seconds of %s : %.2f\n", jsonData.HardwareID, avgTemperature)

	}
}
