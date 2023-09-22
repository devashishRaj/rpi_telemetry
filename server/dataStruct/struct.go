package datastruct

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
	ProcesN     int64   `json:"nprocs"`
}
