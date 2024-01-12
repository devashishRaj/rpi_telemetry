package postgresDB

import (
	"context"
	"database/sql"

	dataStruct "devashishRaj/rpi_telemetry/server/dataStruct"
	handle "devashishRaj/rpi_telemetry/server/handleError"
)

func AlertTemp(jsonData dataStruct.MetricsBatch) {
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
	err := G_dbpool.QueryRow(context.Background(), query, jsonData.MacAddr).Scan(&avgTemperature)
	handle.CheckError("Error executing avgTemp query", err)

	if avgTemperature.Valid {
		if avgTemperature.Float64 > 44 {
			_, err := G_dbpool.Exec(context.Background(), `
			INSERT INTO telemetry.rpi_temp_alert (macaddress , name , value ,timestamp)
			VALUES ($1, $2, $3)`,
				jsonData.MacAddr, deviceMetrics[0].Name,
				deviceMetrics[0].Value,
				deviceMetrics[0].TimeStamp)

			handle.CheckError("error inserting avgTemp above threshold in DB", err)
		}

	}

}
