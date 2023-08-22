package scraprpi

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

type SystemInfo struct {
	HardwareID  string  `json:"HardwareID"`
	CPUuserLoad float64 `json:"CPUuserLoad"`
	CPUidle     float64 `json:"CPUidle"`
	TotalMemory int64   `json:"TotalMemory"`
	FreeMemory  int64   `json:"FreeMemory"`
	IP          string  `json:"IP"`
	Temperature float64 `json:"Temperature"`
	TimeStamp   string  `json:"TimeStamp"`
}

// global var to store update values
var G_systemInfo SystemInfo

// get raspberry serial nuber , can act as UNIQUE key
func GetRaspberryPiHWID() string {
	cmd := exec.Command("cat", "/sys/firmware/devicetree/base/serial-number")
	output, err := cmd.Output()
	CheckError(err)

	rpiHWID := string(output)
	rpiHWID = strings.TrimRight(rpiHWID, "\u0000")
	return rpiHWID
}

func CalculateCPUUsage(mode string) (float64, error) {
	before, err := cpu.Get()
	CheckError(err)

	time.Sleep(time.Duration(1) * time.Second)

	after, err := cpu.Get()
	CheckError(err)

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

func GetMemoryValue(mode string) (float64, error) {
	memory, err := memory.Get()
	CheckError(err)

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

func GetPublicIPAddress() (string, error) {
	cmd := exec.Command("hostname", "-I")
	output, err := cmd.Output()
	CheckError(err)

	ipAddress := strings.TrimSpace(string(output))
	return ipAddress, nil
}

func GetInternalTemperature() float64 {
	cmd := exec.Command("vcgencmd", "measure_temp")
	output, err := cmd.Output()
	CheckError(err)

	temperature := strings.TrimSpace(string(output))
	// Define the regular expression pattern
	pattern := `\d+\.\d+`

	// Compile the regular expression
	regExp := regexp.MustCompile(pattern)

	// Find the first match in the input string
	match := regExp.FindString(temperature)
	numericValue, err := strconv.ParseFloat(match, 64)
	CheckError(err)
	return numericValue
}

// to store values into struct
func StartScraping() SystemInfo {

	// Scrape /proc/stat or gather your data here
	currentTime := time.Now()
	rpiTimeStamp := currentTime.Format("2006-01-02 15:04:05")

	cpuUserUsage, err := CalculateCPUUsage("user")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}

	cpuIdle, err := CalculateCPUUsage("idle")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}

	memTotal, err := GetMemoryValue("total")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}

	memFree, err := GetMemoryValue("free")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}

	publicIP, err := GetPublicIPAddress()
	if err != nil {
		log.Fatal(err)
	}

	internalTemp := GetInternalTemperature()

	// Create a struct to hold the data
	systemInfo := SystemInfo{
		HardwareID:  GetRaspberryPiHWID(),
		CPUuserLoad: cpuUserUsage,
		CPUidle:     cpuIdle,
		TotalMemory: int64(memTotal),
		FreeMemory:  int64(memFree),
		IP:          publicIP,
		Temperature: internalTemp,
		TimeStamp:   rpiTimeStamp,
	}

	return systemInfo

}
