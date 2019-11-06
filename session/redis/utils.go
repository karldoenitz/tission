package redis

import (
	"crypto/md5"
	"fmt"
	"net"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func getGoroutineId() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

// getLocalIP returns the non loopback local IP of the host
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func toMd5(original string) string {
	bytes := []byte(original)
	has := md5.Sum(bytes)
	return fmt.Sprintf("%X", has)
}

func getSessionId() string {
	ip := getLocalIP()
	goId := getGoroutineId()
	timeStamp := time.Now().Local().Unix()
	origin := fmt.Sprintf("%s-%d-%d", ip, goId, timeStamp)
	return toMd5(origin)
}
