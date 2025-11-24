package braza

import (
	"net"
	"reflect"
	"runtime"
	"strings"
)

func getFunctionName(i interface{}) string {
	splitName := strings.Split(
		runtime.FuncForPC(
			reflect.ValueOf(i).Pointer(),
		).Name(), ".",
	)
	return splitName[len(splitName)-1]
}

// Get preferred outbound ip of this machine
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "0.0.0.0"
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}
