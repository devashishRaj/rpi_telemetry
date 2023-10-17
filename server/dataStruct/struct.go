package dataStruct

import "time"

type SystemMetrics struct {
	Name      string
	Value     float64
	TimeStamp time.Time
}

// SystemInfo represents the system information.
type SystemInfo struct {
	MacAddress string `json:"MacAddress"`
	PrivateIP  string `json:"privateIP"`
	PublicIP   string `json:"publicIP"`
	Hostname   string `json:"hostname"`
	OsType     string `json:"ostype"`
}

type MetricsBatch struct {
	MacAddr string
	Metrics []SystemMetrics
}

// type SystemMetrics struct {
// 	CPUuserLoad float64   `json:"CPUuserLoad"`
// 	TotalMemory int64     `json:"TotalMemory"`
// 	FreeMemory  int64     `json:"FreeMemory"`
// 	Temperature float64   `json:"Temperature"`
// 	TimeStamp   time.Time `json:"TimeStamp"`
// 	ProcesN     int64     `json:"nprocs"`
// }
