package senddata

import (
	"bytes"
	handlerror "devashishRaj/rpi_telemetry/client/Handlerror"
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
	handlerror.CheckError("unable to read config , func: ReadConfig", err)
}

func HttpPost(input interface{}, dataflag string) {
	ReadConfig()

	jsonData, err := json.Marshal(input)
	handlerror.CheckError("error during marshal", err)
	var response Response
	var URL string
	if dataflag == "metrics" {
		//URL = "http://10.147.19.40:8080/tele/metrics"
		URL = viper.GetString("baseURL") + "/tele/metrics"
	} else if dataflag == "sysinfo" {
		//URL = "http://10.147.19.40:8080/tele/sysinfo"
		URL = viper.GetString("baseURL") + "/tele/sysinfo"
	}

	resp, err := http.Post(URL, "application/json", bytes.NewBuffer(jsonData))
	handlerror.CheckError("HTTP POST error", err)
	defer resp.Body.Close()
	log.Println("Response Status:", resp.Status)

	responseBody, err := io.ReadAll(resp.Body)
	handlerror.CheckError("Error handling response body", err)

	err = json.Unmarshal(responseBody, &response)
	handlerror.CheckError("error during unmarshal of responseBody", err)

	fmt.Println("Received message:", response.Message)

}
