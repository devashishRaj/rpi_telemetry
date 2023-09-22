package postgresDB

import (
	"database/sql"
	"fmt"
	dataStruct "server/dataStruct"

	_ "github.com/lib/pq"
)

func CheckPrimaryKey(jsonData dataStruct.SystemInfo) {
	db = ConnectDB()
	query := (`
		SELECT EXISTS (SELECT hardwareid , privateIP ,  publicIP , hostname , ostype , totalmemory
			  FROM telemetry.devices 
			WHERE hardwareid = $1 AND
			privateip = $2 AND
			publicIP = $3 AND
			hostname  = $4 AND 
			ostype = $5 AND
			totalmemory = $6);
		`)
	var isPresent bool
	err := db.QueryRow(query, jsonData.HardwareID, jsonData.PrivateIP,
		jsonData.PublicIP,
		jsonData.Hostname, jsonData.OsType,
		jsonData.TotalMemory).Scan(&isPresent)

	if err != nil {
		fmt.Println("Query error in CheckPrimary")
		fmt.Println(err)

	} else if isPresent {
		(InsertInDB(jsonData, db))

	} else {
		(insertHwID(jsonData, db))
	}

}

func insertHwID(jsonData dataStruct.SystemInfo, db *sql.DB) {
	_, err := db.Exec(`
	UPDATE telemetry.devices
	SET privateip = $1, publicIP = $2, hostname = $3, ostype = $4, totalmemory = $5
	WHERE HardwareID = $6`,
		jsonData.PrivateIP, jsonData.PublicIP,
		jsonData.Hostname, jsonData.OsType,
		jsonData.TotalMemory, jsonData.HardwareID)

	if err != nil {
		fmt.Println("Query error in Insert HardWareID")
		fmt.Println(err)

	} else {
		fmt.Println("New HardwareID inserted!")
		InsertInDB(jsonData, db)
	}

}

func InsertInDB(jsonData dataStruct.SystemInfo, db *sql.DB) {
	_, err := db.Exec(`
	INSERT INTO telemetry.rpi4b_metrics (HardwareID, CPUuserLoad,  MemoryUsage,  
								Temperature, TotalProcesses , TimeStamp)
	VALUES ($1, $2, $3, $4, $5, $6)`,
		jsonData.HardwareID, jsonData.CPUuserLoad,
		jsonData.TotalMemory-jsonData.FreeMemory,
		jsonData.Temperature, jsonData.ProcesN, jsonData.TimeStamp)
	if err != nil {

		fmt.Println("error in InsertInDB")
		fmt.Println(err)

	} else {
		fmt.Println("Data inserted successfully!")
		AlertTemp(jsonData, db)
	}
}
