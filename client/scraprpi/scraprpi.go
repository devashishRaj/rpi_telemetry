package scraprpi

import (
	"fmt"
	"log"
	"math"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	datastruct "devashishRaj/rpi_telemetry/server/dataStruct"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
)

func isNil(context string, v interface{}) {
	if v == nil {
		fmt.Printf("Error: %s\n", context)
		log.Fatal("null value ")
	}
}

func CheckError(context string, err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", context)
		log.Fatalln(err)
	}
}

func ReutrnMinus(context string, err error) {

	fmt.Printf("Error: %s\n", context)
}

func GetmacAddr() string {
	var cmdStr string
	os := runtime.GOOS
	if os == "darwin" {
		cmdStr = "ifconfig en0 | grep 'ether' | awk '{print $2}'"
		// for zti refer following command
		// "ifconfig | awk '/^feth[0-9]+/{getline; while (getline) { if ($1 == "ether") { print $2; break } } }'"
	} else {
		cmdStr = "ifconfig eth0 | grep 'ether' | awk '{print $2}'"
	}
	cmd := exec.Command("sh", "-c", cmdStr)
	output, err := cmd.Output()
	CheckError("error in macaddress", err)
	isNil("error in macaddress", output)

	macaddress := string(output)
	// use this trim when you get cmd line prompt on same line as output
	macaddress = strings.TrimRight(macaddress, "\u0000")
	macaddress = strings.TrimSpace(string(macaddress))
	return macaddress

}
func Gethostnmae() string {
	cmd := exec.Command("hostname")
	output, err := cmd.Output()
	CheckError("hostname", err)
	isNil("error in host", output)

	hostname := strings.TrimSpace(string(output))
	return hostname
}

func CalculateCPUUsage(mode string, precision int) float64 {
	before, err := cpu.Get()
	CheckError("CPU", err)
	isNil("error in CPU", before)

	time.Sleep(time.Duration(1) * time.Second)

	after, err := cpu.Get()
	CheckError("CPU", err)
	isNil("error in CPU", after)

	total := float64(after.Total - before.Total)

	switch mode {
	case "user":
		cpuUsr := (float64(after.User-before.User) / total * 100)
		isNil("error in CPU", cpuUsr)
		return math.Round(cpuUsr*math.Pow(10, float64(precision))) / math.Pow(10, float64(precision))
	case "idle":
		cpuIdle := (float64(after.Idle-before.Idle) / total * 100)
		isNil("error in CPU", cpuIdle)
		return math.Round(cpuIdle*math.Pow(10, float64(precision))) / math.Pow(10, float64(precision))
	default:
		return 0.0
	}
}

func GetMemoryValue(mode string) float64 {
	memory, err := memory.Get()
	CheckError("memory", err)
	isNil("error in memory", memory)

	switch mode {
	case "total":
		memTotal := float64(memory.Total) / 1000000
		isNil("error in memory", memTotal)
		return memTotal
	case "free":
		memFree := float64(memory.Free) / 1000000
		isNil("error in memory", memFree)
		return memFree
	default:
		return 0.0
	}
}

func GetPrivateIP() string {
	os := runtime.GOOS
	var cmdStr string
	if os == "darwin" {
		cmdStr = "ifconfig | awk '/^feth[0-9]+/{getline; while (getline) { if ($1 == \"inet\") { print $2; break } } }'"
	} else {
		cmdStr = "ifconfig | awk '/^.*zti[^:]*:/ {while (getline) { if ($1 == \"inet\") { print $2; break } } }'"
	}
	cmd := exec.Command("sh", "-c", cmdStr)
	output, err := cmd.Output()
	CheckError("privateIP", err)
	isNil("error in privateIP", output)

	ipAddress := strings.TrimSpace(string(output))
	return ipAddress
}

func GetPublicIP() string {
	cmd := exec.Command("curl", "ipinfo.io/ip")
	output, err := cmd.Output()
	CheckError("publicIP", err)
	isNil("error in publicIP", output)

	pubAddress := strings.TrimSpace(string(output))
	pubAddress = strings.TrimRight(pubAddress, "\u0000")
	fmt.Println(pubAddress)
	return pubAddress
}

func GetInternalTemperature() float64 {
	cmd := exec.Command("vcgencmd", "measure_temp")
	output, err := cmd.Output()
	if err != nil {
		ReutrnMinus("temp", err)
		return -1
	}

	temperature := strings.TrimSpace(string(output))
	// Define the regular expression pattern
	pattern := `\d+\.\d+`

	// Compile the regular expression
	regExp := regexp.MustCompile(pattern)

	// Find the first match in the input string
	match := regExp.FindString(temperature)
	numericValue, err := strconv.ParseFloat(match, 64)
	CheckError("temp string parse", err)
	return numericValue
}

func getOStype() string {
	os := runtime.GOOS
	return os
}

func TotalProcesses() int64 {
	cmdStr := "ps aux | wc -l"
	cmd := exec.Command("sh", "-c", cmdStr)

	// Capture the command output
	output, err := cmd.Output()
	if err != nil {
		ReutrnMinus("TotalProcesses", err)
		return -1
	}
	trimmedOutput := strings.TrimSpace(string(output))
	// Convert the output to an integer
	count, err := strconv.ParseInt(trimmedOutput, 10, 64)
	CheckError("n processes , string parse", err)

	return count

}

// to store values into struct
func ScrapMetrics() datastruct.MetricsBatch {

	MetricsData := datastruct.MetricsBatch{
		MacAddr: GetmacAddr(),
		Metrics: []datastruct.SystemMetrics{
			{
				CPUuserLoad: CalculateCPUUsage("user", 2),
				TotalMemory: int64(GetMemoryValue("total")),
				FreeMemory:  int64(GetMemoryValue("free")),
				Temperature: GetInternalTemperature(),
				TimeStamp:   time.Now().UTC(),
				ProcesN:     TotalProcesses(),
			},
			// Add more SystemMetrics as needed
		},
	}

	return MetricsData

}
func ScrapSysInfo() datastruct.SystemInfo {
	sysinfo := datastruct.SystemInfo{
		MacAddress: GetmacAddr(),
		PrivateIP:  GetPrivateIP(),
		PublicIP:   GetPublicIP(),
		Hostname:   Gethostnmae(),
		OsType:     getOStype(),
	}
	return sysinfo
}
