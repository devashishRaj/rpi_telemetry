package postgresDB

import (
	"database/sql"
	"fmt"
	"log"
	dataStruct "server/dataStruct"
)

func AlertTemp(jsonData dataStruct.SystemInfo, db *sql.DB) {

	// // Set custom prefix and flags to differentiate alerts
	// log.SetPrefix("[ALERT] ")

	// //The log.SetFlags function is used to configure the log message format to include the date,
	// //time, and source file information.

	// log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	query :=
		`
		SELECT AVG(temperature) AS avg_temperature
		FROM telemetry.rpi4b_metrics
		WHERE hardwareID = $1
		AND timestamp >= NOW() - INTERVAL '60 seconds';
		`
	var avgTemperature sql.NullFloat64
	err := db.QueryRow(query, jsonData.HardwareID).Scan(&avgTemperature)
	if err != nil {

		fmt.Println("calc query error in AvgTemp")
		fmt.Println(err)

	}

	if avgTemperature.Valid {
		if avgTemperature.Float64 > 44.5 {
			_, err := db.Exec(`
			INSERT INTO telemetry.rpi_temp_alert (HardwareID, CPUuserLoad, MemoryUse , privateIP ,
				Temperature, TimeStamp)
			VALUES ($1, $2, $3, $4, $5, $6)`,
				jsonData.HardwareID, jsonData.CPUuserLoad,
				jsonData.TotalMemory-jsonData.FreeMemory, jsonData.PrivateIP,
				jsonData.Temperature, jsonData.TimeStamp)

			if err != nil {

				fmt.Println("error in InsertInDB , isnertoin query")
				fmt.Println(err)

			} else {
				fmt.Println("Data inserted successfully!")

				log.Printf("Average temperature within the last 30 seconds of %s : %.2f\n",
					jsonData.HardwareID, avgTemperature.Float64)
			}
		}

	}
}
