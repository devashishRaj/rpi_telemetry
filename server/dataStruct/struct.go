package datastruct

type SystemMetrics struct {
	CPUuserLoad float64 `json:"CPUuserLoad"`
	TotalMemory int64   `json:"TotalMemory"`
	FreeMemory  int64   `json:"FreeMemory"`
	Temperature float64 `json:"Temperature"`
	TimeStamp   string  `json:"TimeStamp"`
	ProcesN     int64   `json:"nprocs"`
}

// SystemInfo represents the system information.
type SystemInfo struct {
	MacAddress string `json:"MacAddress"`
	PrivateIP  string `json:"privateIP"`
	PublicIP   string `json:"publicIP"`
	Hostname   string `json:"hostname"`
	OsType     string `json:"ostype"`
	Metrics    SystemMetrics
}
