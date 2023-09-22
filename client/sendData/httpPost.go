package senddata

import (
	"bytes"
	scrapRpi "client/scraprpi"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

func ReadConfig() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("./local/.configs")
	viper.SetConfigName("config") // Register config file name (no extension)
	viper.SetConfigType("json")   // Look for specific type
	err := viper.ReadInConfig()
	viper.WatchConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
	viper.Debug()
}

type Response struct {
	Message string `json:"message"`
}

func SendSystemInfo(sysinfo scrapRpi.SystemInfo) {

	var response Response
	serverUrl := viper.GetString("server_url")
	fmt.Println(serverUrl)
	//log.Fatalln(" serverUrl to string assertion failed")

	// Convert the struct to JSON
	jsonData, err := json.Marshal(sysinfo)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}

	// https://github.com/spf13/viper#getting-values-from-viper

	resp, err := http.Post(serverUrl, "application/json", bytes.NewBuffer(jsonData))
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

func SendInterval() {

	//serverURL := "http://10.147.19.40:8080/rpi"
	for {
		scrapRpi.G_systemInfo = scrapRpi.StartScraping()
		SendSystemInfo(scrapRpi.G_systemInfo)
		time.Sleep(time.Second * 10)
	}
}
