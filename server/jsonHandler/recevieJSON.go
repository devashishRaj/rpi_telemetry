package jsonHandler

import (
	"fmt"
	"net/http"
	jsonStruct "server/dataStruct"
	jsonToDb "server/postgresDB"

	"github.com/gin-gonic/gin"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func ReceiveJSON() {
	var jsonData jsonStruct.SystemInfo

	// err := json.NewDecoder(r.Body).Decode(&jsonData)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// fmt.Printf("Received Info: %+v\n", jsonData)
	// w.WriteHeader(http.StatusOK)
	// // UpdateDB in insert.go

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.POST("/rpi", func(c *gin.Context) {
		if err := c.BindJSON(&jsonData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Printf("Received Info: %+v\n", jsonData)
		c.JSON(http.StatusOK, gin.H{"message": "JSON data received successfully"})
		DBerr := jsonToDb.InsertInDB(jsonData)
		CheckError(DBerr)
		jsonToDb.CloseDB()
	})

	r.Run(":8080")
}
