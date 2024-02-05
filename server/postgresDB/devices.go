package postgresDB

import (
	"context"
	dataStruct "github.com/devashishRaj/rpi_telemetry/server/dataStruct"
	handle "github.com/devashishRaj/rpi_telemetry/server/handleError"

	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CheckDevicesDB(db *pgxpool.Pool, jsonData dataStruct.SystemInfo) {
	query := (`
		SELECT EXISTS (SELECT macaddress 
			  FROM telemetry.devices 
			WHERE  macaddress = $1 );
		`)
	var isPresent bool
	err := db.QueryRow(context.Background(), query, jsonData.MacAddress).Scan(&isPresent)
	handle.CheckError("Query error in CheckDevicesDB ", err)

	if !isPresent {
		InsertDevice(db, jsonData)

	} else {
		CheckPrimaryKey(db, jsonData)
	}
}

func InsertDevice(db *pgxpool.Pool, jsonData dataStruct.SystemInfo) {
	_, err := db.Exec(context.Background(), `
		INSERT INTO telemetry.devices (MacAddress , privateIP ,  publicIP , hostname , 
		ostype )
		VALUES ($1, $2, $3, $4, $5)`,
		jsonData.MacAddress, jsonData.PrivateIP, jsonData.PublicIP,
		jsonData.Hostname, jsonData.OsType)
	handle.CheckError("Query error when insert new device , func: CheckDevicesDB", err)

}

func CheckPrimaryKey(db *pgxpool.Pool, jsonData dataStruct.SystemInfo) {
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
	err := db.QueryRow(context.Background(), query, jsonData.MacAddress, jsonData.PrivateIP,
		jsonData.PublicIP,
		jsonData.Hostname, jsonData.OsType).Scan(&isOutdated)

	handle.CheckError("Error in Check Primary key ", err)
	if !isOutdated {
		updateDeviceInfo(db, jsonData)
	}

}

func updateDeviceInfo(db *pgxpool.Pool, jsonData dataStruct.SystemInfo) {

	_, err := db.Exec(context.Background(),
		`UPDATE telemetry.devices
		SET privateip = $1, publicip = $2, hostname = $3, ostype = $4
		WHERE MacAddress = $5`,
		jsonData.PrivateIP, jsonData.PublicIP,
		jsonData.Hostname, jsonData.OsType, jsonData.MacAddress)
	handle.CheckError("Query error in Update device info", err)

}
