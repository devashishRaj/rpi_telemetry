package postgresDB

import (
	//"bytes"

	"fmt"

	//"io"

	//"time"

	_ "github.com/lib/pq"
)

func InsertInDB(jsonData SystemInfo) error {
	db = ConnectDB()
	_, err := db.Exec(`
	INSERT INTO telemetry.rpib (HardwareID, CPUuserLoad, CPUidle, TotalMemory, FreeMemory, IP, 
								Temperature, TimeStamp)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		jsonData.HardwareID, jsonData.CPUuserLoad, jsonData.CPUidle,
		jsonData.TotalMemory, jsonData.FreeMemory, jsonData.IP,
		jsonData.Temperature, jsonData.TimeStamp)

	CheckError(err)
	fmt.Println("Data inserted successfully!")
	AlertTemp()
	return nil
}
