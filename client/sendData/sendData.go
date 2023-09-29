package senddata

import (
	"bytes"
	scrapRpi "devashishRaj/rpi_telemetry/client/scraprpi"
	datastruct "devashishRaj/rpi_telemetry/server/dataStruct"
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

func StructToJSON(input interface{}) {
	jsonData, err := json.Marshal(input)
	if err != nil {
		log.Fatalln(err)
	}
	httpPost(bytes.NewBuffer(jsonData))
}

func httpPost(jsonData *bytes.Buffer) {
	var serverURL string
	//ReadConfig()
	//serverURL = viper.GetString("server")
	var response Response
	serverURL = "http://10.147.19.40:8080/rpi"
	resp, err := http.Post(serverURL, "application/json", jsonData)
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

func SendData(sysinfo datastruct.SystemInfo) {

	// https://github.com/spf13/viper#getting-values-from-viper
	// serverUrl := viper.GetString("server_url")
	// fmt.Println(serverUrl)
	//log.Fatalln(" serverUrl to string assertion failed")

	StructToJSON(sysinfo)
}

func SendInterval() {

	for {
		//scrapRpi.G_systemInfo = scrapRpi.StartScraping()
		SendData(scrapRpi.StartScraping())
		time.Sleep(time.Second * 10)
	}
}
