package jsonHandler

import (
	dataStruct "devashishRaj/rpi_telemetry/server/dataStruct"
	postgresDB "devashishRaj/rpi_telemetry/server/postgresDB"

	//"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var sysinfoMutex sync.Mutex

func ReceiveJSON() {
	var metricsData dataStruct.MetricsBatch
	var sysinfoData dataStruct.SystemInfo

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.POST("/tele/metrics", func(c *gin.Context) {
		// Todo(check plugin for fuzzy todo tag) add function to check
		if err := c.BindJSON(&metricsData); err != nil {
			log.Println()
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Println(metricsData)
			return
		} else {
			//fmt.Printf("Received Metrics: %+v\n", metricsData)
			c.JSON(http.StatusOK, gin.H{"message": "Metrics data received successfully"})
			postgresDB.InsertInDB(metricsData)
		}
	})

	r.POST("/tele/sysinfo", func(c *gin.Context) {
		sysinfoMutex.Lock()
		defer sysinfoMutex.Unlock()
		if err := c.BindJSON(&sysinfoData); err != nil {
			log.Println(sysinfoData)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			//fmt.Printf("Received System Info: %+v\n", sysinfoData)
			c.JSON(http.StatusOK, gin.H{"message": "System Info data received successfully"})
			postgresDB.CheckDevicesDB(sysinfoData)
		}
	})

	err := r.Run(":8080")
	log.Fatalln(err)

}
