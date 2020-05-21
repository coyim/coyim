package net

import (
	"net"
	"os"
	"sync"
)

var onceCheckTails sync.Once
var isTailsVal bool

var onceGetLocalIP sync.Once
var localIPVal string

func getLocalIP() string {
	onceGetLocalIP.Do(func() {
		addrs, _ := net.InterfaceAddrs()
		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					localIPVal = ipnet.IP.String()
					return
				}
			}
		}
	})

	return localIPVal
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func isTails() bool {
	onceCheckTails.Do(func() {
		isTailsVal = fileExists("/etc/amnesia/version")
	})

	return isTailsVal
}
