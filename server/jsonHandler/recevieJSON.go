package jsonHandler

import (
	dataStruct "github.com/devashishRaj/rpi_telemetry/server/dataStruct"
	handle "github.com/devashishRaj/rpi_telemetry/server/handleError"
	postgresDB "github.com/devashishRaj/rpi_telemetry/server/postgresDB"

	//"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

var sysinfoMutex sync.Mutex

func ReceiveJSON() {
	var metricsData dataStruct.MetricsBatch
	var sysinfoData dataStruct.SystemInfo
	var db *pgxpool.Pool = postgresDB.ConnectDB()

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.POST("/tele/metrics", func(c *gin.Context) {

		if err := c.BindJSON(&metricsData); err != nil {
			log.Println()
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			//fmt.Printf("Received Metrics: %+v\n", metricsData)
			c.JSON(http.StatusOK, gin.H{"message": "Metrics data received successfully"})
			postgresDB.InsertInDB(db, metricsData)
		}
	})

	r.POST("/tele/sysinfo", func(c *gin.Context) {
		sysinfoMutex.Lock()
		defer sysinfoMutex.Unlock()
		if err := c.BindJSON(&sysinfoData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			//fmt.Printf("Received System Info: %+v\n", sysinfoData)
			c.JSON(http.StatusOK, gin.H{"message": "System Info data received successfully"})
			postgresDB.CheckDevicesDB(db, sysinfoData)
		}
	})
	defer db.Close()
	err := r.Run(":8080")
	handle.CheckError("Error when conencting to port 8080", err)

}
