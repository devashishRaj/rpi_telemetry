package postgresDB

import (
	"database/sql"
	"fmt"
	"log"
	dataStruct "server/dataStruct"

	_ "github.com/lib/pq"
)

func CheckDevicesDB(jsonData dataStruct.SystemInfo) {
	db = ConnectDB()
	query := (`
		SELECT EXISTS (SELECT macaddress 
			  FROM telemetry.devices 
			WHERE  macaddress = $1 );
		`)
	var isPresent bool
	err := db.QueryRow(query, jsonData.MacAddress).Scan(&isPresent)

	if err != nil {
		fmt.Println("Query error in CheckDevicesDB")
		log.Fatalln(err)

	}
	fmt.Println("isPresent: ", isPresent)
	if isPresent == false {
		_, err := db.Exec(`
		INSERT INTO telemetry.devices (MacAddress , privateIP ,  publicIP , hostname , 
		ostype , totalmemory)
		VALUES ($1, $2, $3, $4, $5, $6)`,
			jsonData.MacAddress, jsonData.PrivateIP, jsonData.PublicIP,
			jsonData.Hostname, jsonData.OsType, jsonData.Metrics.TotalMemory)
		if err != nil {
			fmt.Println("Query error in CheckDevicesDB when  ")
			log.Fatalln(err)

		}
		fmt.Println("New device inserted")
		InsertInDB(jsonData, db)

	} else {
		CheckPrimaryKey(jsonData, db)
	}

}

func CheckPrimaryKey(jsonData dataStruct.SystemInfo, db *sql.DB) {
	query := (`
		SELECT EXISTS (SELECT MacAddress , privateIP ,  publicIP , hostname , ostype , totalmemory
			  FROM telemetry.devices 
			WHERE MacAddress = $1 AND
			privateip = $2 AND
			publicIP = $3 AND
			hostname  = $4 AND 
			ostype = $5 AND
			totalmemory = $6);
		`)
	var isPresent bool
	err := db.QueryRow(query, jsonData.MacAddress, jsonData.PrivateIP,
		jsonData.PublicIP,
		jsonData.Hostname, jsonData.OsType,
		jsonData.Metrics.TotalMemory).Scan(&isPresent)

	if err != nil {
		fmt.Println("Query error in CheckPrimary")
		log.Fatalln(err)
	}

	if isPresent {
		(InsertInDB(jsonData, db))

	} else {
		(updateDeviceInfo(jsonData, db))
	}

}

func updateDeviceInfo(jsonData dataStruct.SystemInfo, db *sql.DB) {

	_, err := db.Exec(`
	UPDATE telemetry.devices
	SET privateip = $1, publicIP = $2, hostname = $3, ostype = $4, totalmemory = $5
	WHERE MacAddress = $6`,
		jsonData.PrivateIP, jsonData.PublicIP,
		jsonData.Hostname, jsonData.OsType,
		jsonData.Metrics.TotalMemory, jsonData.MacAddress)

	if err != nil {
		fmt.Println("Query error in Update device info")
		log.Fatalln(err)

	} else {
		fmt.Println("Device info updated")
		InsertInDB(jsonData, db)
	}

}

func InsertInDB(jsonData dataStruct.SystemInfo, db *sql.DB) {
	_, err := db.Exec(`
	INSERT INTO telemetry.rpi4b_metrics (MacAddress, CPUuserLoad,  MemoryUsage,  
								Temperature, TotalProcesses , TimeStamp)
	VALUES ($1, $2, $3, $4, $5, $6)`,
		jsonData.MacAddress, jsonData.Metrics.CPUuserLoad,
		jsonData.Metrics.TotalMemory-jsonData.Metrics.FreeMemory,
		jsonData.Metrics.Temperature, jsonData.Metrics.ProcesN, jsonData.Metrics.TimeStamp)
	if err != nil {

		fmt.Println("error in InsertInDB")
		log.Fatalln(err)

	} else {
		fmt.Println("Data inserted successfully!")
		AlertTemp(jsonData, db)
	}
	CloseDB()
}
