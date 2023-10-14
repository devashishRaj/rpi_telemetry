package main

import (
	sendData "devashishRaj/rpi_telemetry/client/sendData"
)

func main() {

	go sendData.SendMetrics()
	go sendData.SendInfo()
	select {}
}
