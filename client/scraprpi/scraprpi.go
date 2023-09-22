package scraprpi

import (
	"fmt"
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
		fmt.Println(err)
		panic(err)
	}
}

type SystemInfo struct {
	HardwareID  string  `json:"HardwareID"`
	CPUuserLoad float64 `json:"CPUuserLoad"`
	TotalMemory int64   `json:"TotalMemory"`
	FreeMemory  int64   `json:"FreeMemory"`
	PrivateIP   string  `json:"privateIP"`
	PublicIP    string  `json:"publicIP"`
	Temperature float64 `json:"Temperature"`
	TimeStamp   string  `json:"TimeStamp"`
	Hostname    string  `json:"hostname"`
	OsType      string  `json:"ostype"`
}

// global var to store update values
var G_systemInfo SystemInfo

// get raspberry serial nuber , can act as UNIQUE key
func GetRaspberryPiHWID() string {

	cmd := exec.Command("cat", "/sys/firmware/devicetree/base/serial-number")
	output, err := cmd.Output()
	CheckError(err)

	rpiHWID := string(output)
	// use this trim when you get cmd line prompt on same line as output
	rpiHWID = strings.TrimRight(rpiHWID, "\u0000")

	return rpiHWID

}
func Gethostnmae() string {
	cmd := exec.Command("hostname")
	output, err := cmd.Output()
	CheckError(err)

	hostname := strings.TrimSpace(string(output))
	return hostname
}

func CalculateCPUUsage(mode string) float64 {
	before, err := cpu.Get()
	CheckError(err)

	time.Sleep(time.Duration(1) * time.Second)

	after, err := cpu.Get()
	CheckError(err)

	total := float64(after.Total - before.Total)

	switch mode {
	case "user":
		cpuUsr := (float64(after.User-before.User) / total * 100)
		return cpuUsr
	case "idle":
		cpuIdle := (float64(after.Idle-before.Idle) / total * 100)
		return cpuIdle
	default:
		return 0.0
	}
}

func GetMemoryValue(mode string) float64 {
	memory, err := memory.Get()
	CheckError(err)

	switch mode {
	case "total":
		memTotal := float64(memory.Total) / 1000000
		return memTotal
	case "free":
		memFree := float64(memory.Free) / 1000000
		return memFree
	default:
		return 0.0
	}
}

func GetPrivateIP() string {

	cmd := exec.Command("hostname", "-I")
	output, err := cmd.Output()
	CheckError(err)

	ipAddress := strings.TrimSpace(string(output))
	return ipAddress
}

func GetPublicIP() string {
	cmd := exec.Command("curl", "ifconfig.co")
	output, err := cmd.Output()
	CheckError(err)

	pipAddress := strings.TrimSpace(string(output))
	return pipAddress
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

func getOStype() string {
	cmd := exec.Command("lsb_release", "-ds")
	output, err := cmd.Output()
	CheckError(err)

	return strings.TrimSpace(string(output))
}

// to store values into struct
func StartScraping() SystemInfo {

	currentTime := time.Now()

	// Create a struct to hold the data
	systemInfo := SystemInfo{
		HardwareID:  GetRaspberryPiHWID(),
		CPUuserLoad: CalculateCPUUsage("user"),
		TotalMemory: int64(GetMemoryValue("total")),
		FreeMemory:  int64(GetMemoryValue("free")),
		PrivateIP:   GetPrivateIP(),
		PublicIP:    GetPublicIP(),
		Temperature: GetInternalTemperature(),
		TimeStamp:   currentTime.Format("2006-01-02 15:04:05"),
		Hostname:    Gethostnmae(),
		OsType:      getOStype(),
	}

	return systemInfo

}
