package senddata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func HttpPost(input interface{}, dataflag string) {

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
