package main

import (
	"fmt"
	"strings"

	//"net"
	"encoding/json"

	"os"
	"time"

	"log"
	"net/http"
	"os/exec"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
)

type SystemInfo struct {
	HardwareID  string  `json:"HardwareID"`
	CPUuserLoad float64 `json:"CPUuserLoad"`
	CPUidle     float64 `json:"CPUidle"`
	TotalMemory int64   `json:"TotalMemory"`
	FreeMemory  int64   `json:"FreeMemory"`
	IP          string  `json:"IP"`
	Temperature string  `json:"Temperature"`
	TimeStamp   string  `json:"TimeStamp"`
}

// global var to store update values
var g_systemInfo []byte

// get raspberry serial nuber , can act as UNIQUE key
func getRaspberryPiHWID() string {
	cmd := exec.Command("cat", "/sys/firmware/devicetree/base/serial-number")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	rpiHWID := string(output)
	rpiHWID = strings.TrimRight(rpiHWID, "\u0000")
	return rpiHWID
}

func calculateCPUUsage(mode string) (float64, error) {
	before, err := cpu.Get()
	if err != nil {
		return 0.0, err
	}

	time.Sleep(time.Duration(1) * time.Second)

	after, err := cpu.Get()
	if err != nil {
		return 0.0, err
	}

	total := float64(after.Total - before.Total)

	switch mode {
	case "user":
		cpuUsr := (float64(after.User-before.User) / total * 100)
		return cpuUsr, nil
	case "idle":
		cpuIdle := (float64(after.Idle-before.Idle) / total * 100)
		return cpuIdle, nil
	default:
		return 0.0, fmt.Errorf("invalid mode: %s", mode)
	}
}

func getMemoryValue(mode string) (float64, error) {
	memory, err := memory.Get()
	if err != nil {
		return 0.0, err
	}

	switch mode {
	case "total":
		memTotal := float64(memory.Total) / 1000000
		return memTotal, nil
	case "free":
		memFree := float64(memory.Free) / 1000000
		return memFree, nil
	default:
		return 0.0, fmt.Errorf("invalid mode: %s", mode)
	}
}

func getPublicIPAddress() (string, error) {
	cmd := exec.Command("hostname", "-I")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	ipAddress := strings.TrimSpace(string(output))
	return ipAddress, nil
}

func getInternalTemperature() (string, error) {
	cmd := exec.Command("vcgencmd", "measure_temp")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	temperature := strings.TrimSpace(string(output))
	return temperature, nil
}

// to store values into struct
func startScrapingAndSending() {
	for {
		// Scrape /proc/stat or gather your data here
		currentTime := time.Now()

		cpuUserUsage, err := calculateCPUUsage("user")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}

		cpuIdle, err := calculateCPUUsage("idle")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}

		memTotal, err := getMemoryValue("total")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}

		memFree, err := getMemoryValue("free")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}

		publicIP, err := getPublicIPAddress()
		if err != nil {
			log.Fatal(err)
		}

		internalTemp, err := getInternalTemperature()
		if err != nil {
			log.Fatal(err)
		}
		rpiTimeStamp := currentTime.Format("2006-01-02 15:04:05")

		// Create a struct to hold the data
		systemInfo := SystemInfo{
			HardwareID:  getRaspberryPiHWID(),
			CPUuserLoad: cpuUserUsage,
			CPUidle:     cpuIdle,
			TotalMemory: int64(memTotal),
			FreeMemory:  int64(memFree),
			IP:          publicIP,
			Temperature: internalTemp,
			TimeStamp:   rpiTimeStamp,
		}

		// Convert the struct to JSON
		jsonData, err := json.Marshal(systemInfo)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			continue
		}
		g_systemInfo = jsonData
		fmt.Println(string(jsonData))

		// Wait for some time before the next iteration
		time.Sleep(time.Second * 10) // Adjust the interval as needed
	}
}
func serveSystemInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(g_systemInfo)
}

func main() {

	go startScrapingAndSending() // to update values at set interval

	http.HandleFunc("/send-json", serveSystemInfo)
	http.ListenAndServe(":8080", nil)

}
