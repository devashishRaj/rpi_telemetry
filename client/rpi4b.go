package main

import (
	scrapRpi "client/scraprpi"
	send "client/sendData"
	"fmt"

	"time"
)

func main() {

	serverURL := "http://192.168.1.3:8080/rpi"
	for {
		scrapRpi.G_systemInfo = scrapRpi.StartScraping()
		sendErr := send.SendSystemInfo(serverURL, scrapRpi.G_systemInfo)
		if sendErr != nil {
			fmt.Println("Error sending CPU info:", sendErr)
			return
		}

		fmt.Println("CPU info sent to server")

		time.Sleep(time.Second * 10)
	}

}
