package senddata

import (
	"bytes"
	handlerror "github.com/devashishRaj/rpi_telemetry/client/Handlerror"
	"encoding/json"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

type Response struct {
	Message string `json:"message"`
}

func ReadConfig() {
	viper.AddConfigPath("./local")
	viper.SetConfigName("config") // Register config file name (no extension)
	viper.SetConfigType("json")   // Look for specific type
	err := viper.ReadInConfig()
	handlerror.CheckError("unable to read config , func: ReadConfig", err)
}

func UrlHandler(input interface{}, dataflag string) {
	ReadConfig()

	jsonData, err := json.Marshal(input)
	handlerror.CheckError("error during marshal", err)

	var URL string
	switch dataflag {
	case "metrics":
		URL = viper.GetString("baseURL") + "/tele/metrics"
		HttpPost(URL, jsonData)
	case "sysinfo":
		URL = viper.GetString("baseURL") + "/tele/sysinfo"
		HttpPost(URL, jsonData)
	default:
		log.Fatalln("invalid value in dataFlag")
	}
}

func HttpPost(URL string, jsonData []byte) {
	resp, err := http.Post(URL, "application/json", bytes.NewBuffer(jsonData))
	handlerror.CheckError("HTTP POST error", err)
	defer resp.Body.Close()
	log.Println("Response Status:", resp.Status)
}
