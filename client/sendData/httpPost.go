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
)

type Response struct {
	Message string `json:"message"`
}

func SendSystemInfo(serverURL string, sysinfo scrapRpi.SystemInfo) error {
	var response Response
	// Convert the struct to JSON
	jsonData, err := json.Marshal(sysinfo)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}

	resp, err := http.Post(serverURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("HTTP POST error: %w", err)
	}

	log.Println("Response Status:", resp.Status)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading response body: %w", err)
	}

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	fmt.Println("Received message:", response.Message)
	defer resp.Body.Close()

	return nil
}

func SendInterval() {

	serverURL := "http://192.168.1.3:8080/rpi"
	for {
		scrapRpi.G_systemInfo = scrapRpi.StartScraping()
		sendErr := SendSystemInfo(serverURL, scrapRpi.G_systemInfo)
		if sendErr != nil {
			fmt.Println("Error sending CPU info:", sendErr)
			return
		}

		time.Sleep(time.Second * 10)
	}
}
