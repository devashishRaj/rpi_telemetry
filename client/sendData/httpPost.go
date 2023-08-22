package senddata

import (
	"bytes"
	scrapRpi "client/scraprpi"
	"fmt"

	"encoding/json"
	"log"
	"net/http"
)

func SendSystemInfo(serverURL string, sysinfo scrapRpi.SystemInfo) error {
	// Convert the struct to JSON
	jsonData, err := json.Marshal(scrapRpi.G_systemInfo)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}

	resp, err := http.Post(serverURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	log.Println("Response Status:", resp.Status)
	return nil
}
