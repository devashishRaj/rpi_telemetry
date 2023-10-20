package postgresDB

import (
	"context"
	dataStruct "devashishRaj/rpi_telemetry/server/dataStruct"
	"log"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5"
)

func InsertInDB(jsonData dataStruct.MetricsBatch) {
	if len(jsonData.Metrics) == 0 {
		log.Fatal("jsonData is empty")
	}

	var rows [][]interface{}
	for _, metric := range jsonData.Metrics {
		rows = append(rows,
			[]interface{}{jsonData.MacAddr, metric.Name, metric.Value, metric.TimeStamp})
	}
	// _, err := G_dbpool.Exec(context.Background(),
	// 	`COPY telemetry.metrics_new (macaddress, name, value, timestamp) FROM stdin`,
	// 	pgx.CopyFromRows(rows))
	// if err != nil {
	// 	log.Println(rows)
	// 	log.Fatal("Error during bulk insert:", err)
	// }
	copyCount, err := G_dbpool.CopyFrom(
		context.Background(),
		pgx.Identifier{"telemetry", "metrics_new"},
		[]string{"macaddress", "name", "value", "timestamp"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		log.Println(copyCount)
		log.Fatal(err)
	}

	AlertTemp(jsonData)
}
