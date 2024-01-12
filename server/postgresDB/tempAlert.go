package postgresDB

import (
	"context"
	"database/sql"

	dataStruct "devashishRaj/rpi_telemetry/server/dataStruct"
	handle "devashishRaj/rpi_telemetry/server/handleError"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AlertTemp(db *pgxpool.Pool, jsonData dataStruct.MetricsBatch) {
	deviceMetrics := jsonData.Metrics

	query :=
		`
		SELECT AVG(value) AS avg_temperature
		FROM telemetry.metrics_new
		WHERE MacAddress = $1
		AND name = 'temperature'
		AND timestamp >= NOW() - INTERVAL '60 seconds';
		`
	var avgTemperature sql.NullFloat64
	err := db.QueryRow(context.Background(), query, jsonData.MacAddr).Scan(&avgTemperature)
	handle.CheckError("Error executing avgTemp query", err)

	if avgTemperature.Valid {
		if avgTemperature.Float64 > 44 {
			_, err := db.Exec(context.Background(), `
			INSERT INTO telemetry.rpi_temp_alert (macaddress , name , value ,timestamp)
			VALUES ($1, $2, $3)`,
				jsonData.MacAddr, deviceMetrics[0].Name,
				deviceMetrics[0].Value,
				deviceMetrics[0].TimeStamp)

			handle.CheckError("error inserting avgTemp above threshold in DB", err)
		}

	}

}
