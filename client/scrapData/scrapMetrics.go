package scraprpi

import (
	"math"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	handlerror "devashishRaj/rpi_telemetry/client/Handlerror"
	sendData "devashishRaj/rpi_telemetry/client/sendData"
	datastruct "devashishRaj/rpi_telemetry/server/dataStruct"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
)

var metricsMutex sync.Mutex

func CalculateCPUUsage(mode string, precision int) {
	before, err := cpu.Get()
	handlerror.CheckError("CPU", err)
	handlerror.IsNil("error in CPU", before)

	time.Sleep(time.Duration(1) * time.Second)

	after, err := cpu.Get()
	handlerror.CheckError("CPU", err)
	handlerror.IsNil("error in CPU", after)

	total := float64(after.Total - before.Total)

	switch mode {
	case "user":
		cpuUsr := (float64(after.User-before.User) / total * 100)
		handlerror.IsNil("error in CPU", cpuUsr)
		cpuUsr = math.Round(cpuUsr*math.Pow(10, float64(precision))) / math.Pow(10, float64(precision))
		ScrapMetrics("CPUuserLoad", cpuUsr)
	case "idle":
		cpuIdle := (float64(after.Idle-before.Idle) / total * 100)
		handlerror.IsNil("error in CPU", cpuIdle)
		cpuIdle = math.Round(cpuIdle*math.Pow(10, float64(precision))) / math.Pow(10, float64(precision))
		ScrapMetrics("CPUIdle", cpuIdle)
	default:

	}
}

func GetMemoryValue(mode string) {
	memory, err := memory.Get()
	handlerror.CheckError("memory", err)
	handlerror.IsNil("error in memory", memory)

	switch mode {
	case "total":
		memTotal := float64(memory.Total) / 1000000
		handlerror.IsNil("error in memory", memTotal)
		ScrapMetrics("TotalMemory", memTotal)
	case "free":
		memFree := float64(memory.Free) / 1000000
		handlerror.IsNil("error in memory", memFree)
		ScrapMetrics("FreeMemory", memFree)
	default:

	}
}

func GetInternalTemperature() {
	cmd := exec.Command("vcgencmd", "measure_temp")
	output, err := cmd.Output()
	if err != nil {
		handlerror.ReutrnMinus("temp", err)
		ScrapMetrics("Temperature", -1)
	} else {

		temperature := strings.TrimSpace(string(output))
		// Define the regular expression pattern
		pattern := `\d+\.\d+`

		// Compile the regular expression
		regExp := regexp.MustCompile(pattern)

		// Find the first match in the input string
		match := regExp.FindString(temperature)
		numericValue, err := strconv.ParseFloat(match, 64)
		handlerror.CheckError("temp string parse", err)
		ScrapMetrics("Temperature", numericValue)
	}
}

func TotalProcesses() {
	cmdStr := "ps aux | wc -l"
	cmd := exec.Command("sh", "-c", cmdStr)

	// Capture the command output
	output, err := cmd.Output()
	if err != nil {
		handlerror.ReutrnMinus("TotalProcesses", err)
		ScrapMetrics("Temperature", -1)
	} else {
		trimmedOutput := strings.TrimSpace(string(output))
		// Convert the output to an integer
		count, err := strconv.ParseInt(trimmedOutput, 10, 64)
		handlerror.CheckError("n processes , string parse", err)
		ScrapMetrics("Temperature", float64(count))
	}

}

// to store values into struct

func MetricInterval() {
	ticker1 := time.NewTicker(10 * time.Second)
	defer ticker1.Stop()
	for range ticker1.C {
		CalculateCPUUsage("user", 2)
		GetMemoryValue("total")
		GetMemoryValue("free")
	}
	ticker2 := time.NewTicker(30 * time.Second)
	defer ticker1.Stop()
	for range ticker2.C {

		GetInternalTemperature()
		TotalProcesses()
	}

}

func ScrapMetrics(metricsName string, metricValue float64) {
	metricsMutex.Lock()
	defer metricsMutex.Unlock()

	metricsData := datastruct.MetricsBatch{
		MacAddr: GetmacAddr(),
		Metrics: []datastruct.SystemMetrics{
			{
				Name:      metricsName,
				Value:     metricValue,
				TimeStamp: time.Now().UTC(),
			},
		},
	}

	sendData.HttpPost(metricsData, "metrics")

}
