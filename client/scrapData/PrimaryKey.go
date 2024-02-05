package scraprpi

import (
	"os/exec"
	"runtime"
	"strings"

	handlerror "github.com/devashishRaj/rpi_telemetry/client/Handlerror"
)

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
	handlerror.CheckError("error in macaddress", err)
	handlerror.IsNil("error in macaddress", output)

	macaddress := string(output)
	// use this trim when you get cmd line prompt on same line as output
	macaddress = strings.TrimRight(macaddress, "\u0000")
	macaddress = strings.TrimSpace(string(macaddress))
	return macaddress

}
