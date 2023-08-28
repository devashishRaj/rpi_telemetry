package postgresDB

import (
	"fmt"
	"log"
	dataStruct "server/dataStruct"
)

func AlertTemp(jsonData dataStruct.SystemInfo) {

	// // Set custom prefix and flags to differentiate alerts
	// log.SetPrefix("[ALERT] ")

	// //The log.SetFlags function is used to configure the log message format to include the date,
	// //time, and source file information.

	// log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	query :=
		`
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
		db = ConnectDB()
		_, err := db.Exec(`
		INSERT INTO telemetry.rpi_temp_alert (HardwareID, CPUuserLoad, CPUidle, TotalMemory, 
		FreeMemory,IP, Temperature, TimeStamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			jsonData.HardwareID, jsonData.CPUuserLoad, jsonData.CPUidle,
			jsonData.TotalMemory, jsonData.FreeMemory, jsonData.IP,
			jsonData.Temperature, jsonData.TimeStamp)

		CheckError(err)
		fmt.Println("Data inserted successfully!")

		log.Printf("Average temperature within the last 30 seconds of %s : %.2f\n",
			jsonData.HardwareID, avgTemperature)

	}
}
