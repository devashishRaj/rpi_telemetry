package scraprpi

import (
	"os/exec"
	"runtime"
	"strings"
	"time"

	handlerror "github.com/devashishRaj/rpi_telemetry/client/Handlerror"
	sendData "github.com/devashishRaj/rpi_telemetry/client/sendData"
	datastruct "github.com/devashishRaj/rpi_telemetry/server/dataStruct"
)

func Gethostnmae() string {
	cmd := exec.Command("hostname")
	output, err := cmd.Output()
	handlerror.CheckError("hostname", err)
	handlerror.IsNil("error in host", output)

	hostname := strings.TrimSpace(string(output))
	return hostname
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
	handlerror.CheckError("privateIP", err)
	handlerror.IsNil("error in privateIP", output)

	ipAddress := strings.TrimSpace(string(output))
	return ipAddress
}

func GetPublicIP() string {
	cmd := exec.Command("curl", "ipinfo.io/ip")
	output, err := cmd.Output()
	handlerror.CheckError("publicIP", err)
	handlerror.IsNil("error in publicIP", output)

	pubAddress := strings.TrimSpace(string(output))
	pubAddress = strings.TrimRight(pubAddress, "\u0000")
	//fmt.Println(pubAddress)
	return pubAddress
}

func getOStype() string {
	os := runtime.GOOS
	return os
}

// to store values into struct

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

func SendSysInfo() {
	ticker := time.NewTicker(300 * time.Second)
	defer ticker.Stop()

	for range ticker.C {

		sendData.UrlHandler(ScrapSysInfo(), "sysinfo")
	}
}
