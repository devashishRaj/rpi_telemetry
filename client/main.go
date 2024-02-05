package main

import (
	scrapData "github.com/devashishRaj/rpi_telemetry/client/scrapData"
	sendData "github.com/devashishRaj/rpi_telemetry/client/sendData"
)

func main() {
	//for first time as soon program start to check if the device on which this program is running
	// exists in telemetry.devices database or not and thus to add it instantly before metrics
	// is send otherwise will lead to primary key error.
	sendData.UrlHandler(scrapData.ScrapSysInfo(), "sysinfo")

	go scrapData.MetricInterval()
	go scrapData.SendSysInfo()
	select {}
}
