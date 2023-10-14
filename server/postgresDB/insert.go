package postgresDB

import (
	"context"
	dataStruct "devashishRaj/rpi_telemetry/server/dataStruct"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5"
)

func InsertInDB(jsonData dataStruct.MetricsBatch) {
	deviceMetrics := jsonData.Metrics
	_, err := G_dbpool.Exec(context.Background(), `
	INSERT INTO telemetry.rpi4b_metrics (MacAddress, CPUuserLoad,  MemoryUsage,  
								Temperature, TotalProcesses , TimeStamp)
	VALUES ($1, $2, $3, $4, $5, $6)`,
		jsonData.MacAddr, deviceMetrics[0].CPUuserLoad,
		deviceMetrics[0].TotalMemory-deviceMetrics[0].FreeMemory,
		deviceMetrics[0].Temperature, deviceMetrics[0].ProcesN, deviceMetrics[0].TimeStamp)
	if err != nil {

		fmt.Println("error in InsertInDB")
		log.Fatalln(err)

	} else {
		AlertTemp(jsonData)
	}
}
