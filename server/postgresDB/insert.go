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
	if len(deviceMetrics) == 0 {
		log.Fatal("jsondata is emply")
	}

	for _, metric := range deviceMetrics {
		_, err := G_dbpool.Exec(context.Background(), `
            INSERT INTO telemetry.metrics_new (macaddress, name, value, timestamp)
            VALUES ($1, $2, $3, $4)`,
			jsonData.MacAddr, metric.Name, metric.Value, metric.TimeStamp)
		if err != nil {
			fmt.Println("Error inserting metric:", metric.Name)
			log.Fatal(err)
		}
	}

	AlertTemp(jsonData)
}
