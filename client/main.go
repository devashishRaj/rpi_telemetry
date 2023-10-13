package main

import (
	sendData "devashishRaj/rpi_telemetry/client/sendData"
)

func main() {

	sendData.SendMetrics()
	sendData.SendMetrics()

}
