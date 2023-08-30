package postgresDB

import (
	"database/sql"
	dataStruct "server/dataStruct"

	"fmt"

	_ "github.com/lib/pq"
)

func CheckPrimaryKey(jsonData dataStruct.SystemInfo) {
	db = ConnectDB()
	query := (`
		SELECT EXISTS (SELECT hardwareid FROM telemetry.devices WHERE hardwareid = $1);
		`)
	var isPresent bool
	err := db.QueryRow(query, jsonData.HardwareID).Scan(&isPresent)
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
	INSERT INTO telemetry.devices (HardwareID)
	VALUES ($1)`, jsonData.HardwareID)
	fmt.Println("New HardwareID inserted!")
	if err != nil {
		fmt.Println("Query error in InsertHardWareID")

	} else {

		InsertInDB(jsonData, db)
	}

}

func InsertInDB(jsonData dataStruct.SystemInfo, db *sql.DB) {
	_, err := db.Exec(`
	INSERT INTO telemetry.rpi4b_metrics (HardwareID, CPUuserLoad, CPUidle, TotalMemory, FreeMemory, IP, 
								Temperature, TimeStamp)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		jsonData.HardwareID, jsonData.CPUuserLoad, jsonData.CPUidle,
		jsonData.TotalMemory, jsonData.FreeMemory, jsonData.IP,
		jsonData.Temperature, jsonData.TimeStamp)
	if err != nil {

		fmt.Println("error in InsertInDB")
		fmt.Println(err)

	} else {
		fmt.Println("Data inserted successfully!")
		AlertTemp(jsonData, db)
	}
}
