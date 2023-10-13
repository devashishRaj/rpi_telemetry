package postgresDB

import (
	"context"
	dataStruct "devashishRaj/rpi_telemetry/server/dataStruct"
	"fmt"
	"log"

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

	if err != nil {
		fmt.Println("Query error in CheckDevicesDB")
		log.Fatalln(err)

	}
	fmt.Println("isPresent: ", isPresent)
	if !isPresent {
		_, err := G_dbpool.Exec(context.Background(), `
		INSERT INTO telemetry.devices (MacAddress , privateIP ,  publicIP , hostname , 
		ostype )
		VALUES ($1, $2, $3, $4, $5, $6)`,
			jsonData.MacAddress, jsonData.PrivateIP, jsonData.PublicIP,
			jsonData.Hostname, jsonData.OsType)
		if err != nil {
			fmt.Println("Query error in CheckDevicesDB when  ")
			log.Fatalln(err)

		}
		fmt.Println("New device inserted")

	} else {
		CheckPrimaryKey(jsonData)
	}
}

func CheckPrimaryKey(jsonData dataStruct.SystemInfo) {
	query := (`
		SELECT EXISTS (SELECT MacAddress , privateIP ,  publicIP , hostname , ostype , totalmemory
			  FROM telemetry.devices 
			WHERE MacAddress = $1 AND
			privateip = $2 AND
			publicIP = $3 AND
			hostname  = $4 AND 
			ostype = $5 AND);
		`)
	var isPresent bool
	err := G_dbpool.QueryRow(context.Background(), query, jsonData.MacAddress, jsonData.PrivateIP,
		jsonData.PublicIP,
		jsonData.Hostname, jsonData.OsType).Scan(&isPresent)

	if err != nil {
		fmt.Println("Query error in CheckPrimary")
		log.Fatalln(err)
	}
	if !isPresent {
		updateDeviceInfo(jsonData)
	}

}

func updateDeviceInfo(jsonData dataStruct.SystemInfo) {

	_, err := G_dbpool.Exec(context.Background(), `
	UPDATE telemetry.devices
	SET privateip = $1, publicIP = $2, hostname = $3, ostype = $4, 
	WHERE MacAddress = $5`,
		jsonData.PrivateIP, jsonData.PublicIP,
		jsonData.Hostname, jsonData.OsType, jsonData.MacAddress)

	if err != nil {
		fmt.Println("Query error in Update device info")
		log.Fatalln(err)

	}

}
