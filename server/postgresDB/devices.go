package postgresDB

import (
	"context"
	dataStruct "devashishRaj/rpi_telemetry/server/dataStruct"
	handle "devashishRaj/rpi_telemetry/server/handleError"

	_ "github.com/jackc/pgx/v5"
)

func CheckDevicesDB(jsonData dataStruct.SystemInfo) {
	query := (`
		SELECT EXISTS (SELECT macaddress 
			  FROM telemetry.devices 
			WHERE  macaddress = $1 );
		`)
	var isPresent bool
	err := G_dbpool.QueryRow(context.Background(), query, jsonData.MacAddress).Scan(&isPresent)
	handle.CheckError("Query error in CheckDevicesDB ", err)

	if !isPresent {
		_, err := G_dbpool.Exec(context.Background(), `
		INSERT INTO telemetry.devices (MacAddress , privateIP ,  publicIP , hostname , 
		ostype )
		VALUES ($1, $2, $3, $4, $5)`,
			jsonData.MacAddress, jsonData.PrivateIP, jsonData.PublicIP,
			jsonData.Hostname, jsonData.OsType)
		handle.CheckError("Query error when insert new device , func: CheckDevicesDB", err)

	} else {
		CheckPrimaryKey(jsonData)
	}
}

func CheckPrimaryKey(jsonData dataStruct.SystemInfo) {
	query := (`
		SELECT EXISTS (SELECT MacAddress , privateIP ,  publicIP , hostname , ostype 
			FROM telemetry.devices 
			WHERE MacAddress = $1 AND
			privateip = $2 AND
			publicIP = $3 AND
			hostname  = $4 AND 
			ostype = $5 );
		`)
	var isOutdated bool
	err := G_dbpool.QueryRow(context.Background(), query, jsonData.MacAddress, jsonData.PrivateIP,
		jsonData.PublicIP,
		jsonData.Hostname, jsonData.OsType).Scan(&isOutdated)

	handle.CheckError("Error in Check Primary key ", err)
	if !isOutdated {
		updateDeviceInfo(jsonData)
	}

}

func updateDeviceInfo(jsonData dataStruct.SystemInfo) {

	_, err := G_dbpool.Exec(context.Background(),
		`UPDATE telemetry.devices
		SET privateip = $1, publicip = $2, hostname = $3, ostype = $4
		WHERE MacAddress = $5`,
		jsonData.PrivateIP, jsonData.PublicIP,
		jsonData.Hostname, jsonData.OsType, jsonData.MacAddress)
	handle.CheckError("Query error in Update device info", err)

}
