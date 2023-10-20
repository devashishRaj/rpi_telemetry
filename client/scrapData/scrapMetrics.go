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
var accumulatedMetrics []datastruct.SystemMetrics

// CalculateCPUUsage calculates CPU usage based on the given mode and precision.
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
		AccumulateMetrics("CPUuserLoad", cpuUsr)
	case "idle":
		cpuIdle := (float64(after.Idle-before.Idle) / total * 100)
		handlerror.IsNil("error in CPU", cpuIdle)
		cpuIdle = math.Round(cpuIdle*math.Pow(10, float64(precision))) / math.Pow(10, float64(precision))
		AccumulateMetrics("CPUIdle", cpuIdle)
	default:

	}
}

// GetMemoryValue retrieves memory information based on the given mode.
func GetMemoryUsage() {
	// obtained in bytes , divide by 1000 for kb 1000,000 for mb .
	memory, err := memory.Get()
	handlerror.CheckError("memory", err)
	handlerror.IsNil("error in memory", memory)

	memUsed := float64((memory.Used) / 1000000)
	handlerror.IsNil("error in memory", memUsed)

	memFree := float64(memory.Free) / 1000000
	handlerror.IsNil("error in memory", memFree)

	AccumulateMetrics("ramusage", memUsed)

}

// GetInternalTemperature measures the internal temperature and accumulates the result.
func GetInternalTemperature() {
	cmd := exec.Command("vcgencmd", "measure_temp")
	output, err := cmd.Output()
	if err != nil {
		handlerror.ReutrnMinus("temp", err)
		AccumulateMetrics("Temperature", -1)
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
		AccumulateMetrics("Temperature", numericValue)
	}
}

// TotalProcesses counts the total number of processes and accumulates the result.
func TotalProcesses() {
	cmdStr := "ps aux | wc -l"
	cmd := exec.Command("sh", "-c", cmdStr)
	// Capture the command output
	output, err := cmd.Output()
	if err != nil {
		handlerror.ReutrnMinus("TotalProcesses", err)
		AccumulateMetrics("Temperature", -1)
	} else {
		trimmedOutput := strings.TrimSpace(string(output))
		// Convert the output to an integer
		count, err := strconv.ParseInt(trimmedOutput, 10, 64)
		handlerror.CheckError("n processes, string parse", err)
		AccumulateMetrics("TotalProcesses", float64(count))
	}
}

// MetricInterval sets up periodic metric collection.
func MetricInterval() {
	ticker1 := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker1.C {
			CalculateCPUUsage("user", 2)
			GetInternalTemperature()
			SendAccumulatedMetrics()
		}
	}()
	defer ticker1.Stop()
	ticker2 := time.NewTicker(31 * time.Second)
	go func() {
		for range ticker2.C {
			//GetMemoryValue("total")
			GetMemoryUsage()
			TotalProcesses()
			//fmt.Println("Inside ticker 2")
			SendAccumulatedMetrics()
		}
	}()
	defer ticker2.Stop()
	select {}
}

// AccumulateMetrics accumulates a metric in the accumulatedMetrics slice.
func AccumulateMetrics(metricsName string, metricValue float64) {
	metricsMutex.Lock()
	defer metricsMutex.Unlock()

	metric := datastruct.SystemMetrics{
		Name:      metricsName,
		Value:     metricValue,
		TimeStamp: time.Now().UTC(),
	}

	accumulatedMetrics = append(accumulatedMetrics, metric)
}

// SendAccumulatedMetrics sends the accumulated metrics to the server.
func SendAccumulatedMetrics() {
	metricsMutex.Lock()
	defer metricsMutex.Unlock()

	if len(accumulatedMetrics) > 0 {
		metricsData := datastruct.MetricsBatch{
			MacAddr: GetmacAddr(),
			Metrics: accumulatedMetrics,
		}

		sendData.HttpPost(metricsData, "metrics")

		// Clear the accumulated metrics after sending
		accumulatedMetrics = nil
	}
}
