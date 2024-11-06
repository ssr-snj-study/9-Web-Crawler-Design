package internal

import (
	"fmt"
	"net"
	"os"
)

func MyIP() []string {
	myIp := []string{}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("Local IP addresses:")
	for _, addr := range addrs {
		// IPv4 주소만 필터링
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				fmt.Println(ipNet.IP.String())
				myIp = append(myIp, ipNet.IP.String())
			}
		}
	}
	return myIp
}
