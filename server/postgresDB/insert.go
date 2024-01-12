package postgresDB

import (
	"context"
	dataStruct "devashishRaj/rpi_telemetry/server/dataStruct"
	handle "devashishRaj/rpi_telemetry/server/handleError"
	"log"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InsertInDB( db *pgxpool.Pool ,  jsonData dataStruct.MetricsBatch) {
	if len(jsonData.Metrics) == 0 {
		log.Fatal("jsonData is empty")
	}

	// Check if data has been collected for 1 minute
	// If 1 minute has passed, perform the bulk insert

	// Create a 2D slice to hold all the rows from the accumulated data
	var rows [][]interface{}
	for _, metric := range jsonData.Metrics {
		rows = append(rows,
			[]interface{}{jsonData.MacAddr, metric.Name, metric.Value, metric.TimeStamp})
	}

	// Perform the bulk insert as before
	copyCount, err := db.CopyFrom(
		context.Background(),
		pgx.Identifier{"telemetry", "metrics_new"},
		[]string{"macaddress", "name", "value", "timestamp"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		handle.CheckError("Error in bulk insert", err)
		log.Println(copyCount)
		log.Println(rows)
		log.Fatal(err)
	}
	AlertTemp(db , jsonData)

}
