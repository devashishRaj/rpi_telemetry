package senddata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

type Response struct {
	Message string `json:"message"`
}

func ReadConfig() {
	viper.AddConfigPath("./local/.configs")
	viper.SetConfigName("config") // Register config file name (no extension)
	viper.SetConfigType("json")   // Look for specific type
	err := viper.ReadInConfig()
	if err != nil {

		fmt.Println("error in client/sendData/sendData.go at ReadConfig()")
		log.Fatalln(err)
	}
}

func HttpPost(input interface{}, dataflag string) {
	ReadConfig()

	jsonData, err := json.Marshal(input)
	if err != nil {
		log.Fatalln(err)
	}
	var response Response
	var URL string
	if dataflag == "metrics" {

		//URL = "http://10.147.19.40:8080/tele/metrics"
		URL = viper.GetString("systemMetrics")

	} else if dataflag == "sysinfo" {

		//URL = "http://10.147.19.40:8080/tele/sysinfo"
		URL = viper.GetString("systemInfo")

	}

	resp, err := http.Post(URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalln("HTTP POST error:", err)

	}

	log.Println("Response Status:", resp.Status)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Error reading response body: ", err)
	}

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		log.Fatalln("unmarshal ", err)
	}
	fmt.Println("Received message:", response.Message)
	defer resp.Body.Close()

}
