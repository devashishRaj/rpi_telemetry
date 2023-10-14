package postgresDB

import (
	"context"
	"database/sql"

	dataStruct "devashishRaj/rpi_telemetry/server/dataStruct"
	"fmt"
	"log"
)

func AlertTemp(jsonData dataStruct.MetricsBatch) {

	// // Set custom prefix and flags to differentiate alerts
	// log.SetPrefix("[ALERT] ")

	// //The log.SetFlags function is used to configure the log message format to include the date,
	// //time, and source file information.

	// log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	query :=
		`
		SELECT AVG(temperature) AS avg_temperature
		FROM telemetry.rpi4b_metrics
		WHERE MacAddress = $1
		AND timestamp >= NOW() - INTERVAL '60 seconds';
		`
	var avgTemperature sql.NullFloat64
	err := G_dbpool.QueryRow(context.Background(), query, jsonData.MacAddr).Scan(&avgTemperature)
	if err != nil {

		fmt.Println("calc query error in AvgTemp")
		log.Fatalln(err)

	}

	if avgTemperature.Valid {
		if avgTemperature.Float64 > 44 {
			_, err := G_dbpool.Exec(context.Background(), `
			INSERT INTO telemetry.rpi_temp_alert (MacAddress, CPUuserLoad,  MemoryUsage,  
				Temperature, TotalProcesses , TimeStamp)
				VALUES ($1, $2, $3, $4, $5, $6)`,
				jsonData.MacAddr, jsonData.Metrics[0].CPUuserLoad,
				jsonData.Metrics[0].TotalMemory-jsonData.Metrics[0].FreeMemory,
				jsonData.Metrics[0].Temperature, jsonData.Metrics[0].ProcesN, jsonData.Metrics[0].TimeStamp)

			if err != nil {

				fmt.Println("error in InsertInDB , isnertoin query")
				log.Fatalln(err)

			} else {
				log.Printf("Average temperature within the last 30 seconds of %s : %.2f\n",
					jsonData.MacAddr, avgTemperature.Float64)
			}
		}

	} else {
		fmt.Println("avg value is null")
	}

}
