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
	INSERT INTO telemetry.metrics_new (macaddress , name , value ,timestamp)
	VALUES ($1, $2, $3 , $4)`,
		jsonData.MacAddr, deviceMetrics[0].Name,
		deviceMetrics[0].Value,
		deviceMetrics[0].TimeStamp)
	if err != nil {

		fmt.Println("error in InsertInDB")
		log.Fatalln(err)

	} else {
		AlertTemp(jsonData)
	}
}
