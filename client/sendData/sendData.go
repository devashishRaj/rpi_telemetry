package senddata

import (
	"bytes"
	scrapRpi "devashishRaj/rpi_telemetry/client/scraprpi"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	//"github.com/spf13/viper"
)

// func ReadConfig() {
// 	viper.AddConfigPath("$HOME/.config/rpiTele/")
// 	viper.AddConfigPath(".")
// 	viper.SetConfigName("config") // Register config file name (no extension)
// 	viper.SetConfigType("json")   // Look for specific type
// 	err := viper.ReadInConfig()
// 	viper.WatchConfig()
// 	if err != nil {
// 		log.Fatalf("Error reading config file: %s", err)
// 	}
// }

type Response struct {
	Message string `json:"message"`
}

func httpPost(input interface{}, dataflag string) {
	//ReadConfig()
	//serverURL = viper.GetString("server")
	jsonData, err := json.Marshal(input)
	if err != nil {
		log.Fatalln(err)
	}
	var response Response
	var URL string
	if dataflag == "metrics" {
		URL = "http://10.147.19.40:8080/tele/metrics"
	} else if dataflag == "sysinfo" {
		URL = "http://10.147.19.40:8080/tele/sysinfo"
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

func SendMetrics() {

	for {
		//scrapRpi.G_systemInfo = scrapRpi.StartScraping()
		httpPost(scrapRpi.ScrapMetrics(), "metrics")
		time.Sleep(time.Second * 10)
	}
}
func SendInfo() {

	for {
		//scrapRpi.G_systemInfo = scrapRpi.StartScraping()
		httpPost(scrapRpi.ScrapSysInfo(), "sysinfo")
		time.Sleep(time.Second * 30)
	}
}
