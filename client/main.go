package main

import (
	scrapData "devashishRaj/rpi_telemetry/client/scrapData"
)

func main() {

	go scrapData.MetricInterval()
	go scrapData.SendSysInfo()
	select {}
}
