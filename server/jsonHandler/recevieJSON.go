package jsonHandler

import (
	"fmt"
	"net/http"
	dataStruct "server/dataStruct"
	postgresDB "server/postgresDB"

	"github.com/gin-gonic/gin"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func ReceiveJSON() {
	var jsonData dataStruct.SystemInfo

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.POST("/rpi", func(c *gin.Context) {
		if err := c.BindJSON(&jsonData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Printf("Received Info: %+v\n", jsonData)
		c.JSON(http.StatusOK, gin.H{"message": "JSON data received successfully"})
		DBerr := postgresDB.InsertInDB(jsonData)
		CheckError(DBerr)
		postgresDB.CloseDB()
	})

	r.Run(":8080")
}
