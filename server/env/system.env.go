package env

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func GetAll() string {
	envVars := os.Environ()
	for _, envVar := range envVars {
		pair := strings.SplitN(envVar, "=", 2)
		if len(pair) == 2 {
			key := pair[0]
			value := pair[1]
			fmt.Printf("%s: %s\n", key, value)
		}
	}
	return strings.Join(envVars, "\n")
}

func GetLocalIP() (string, error) {
	// Try to get all interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, i := range interfaces {
		// Get addresses associated with each interface
		addrs, err := i.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			// Check if the address is an IP address
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// Skip loopback addresses
			if ip == nil || ip.IsLoopback() {
				continue
			}

			// Only consider IPv4 addresses
			if ip.To4() != nil {
				return ip.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no local IP address found")
}

func GetPwd() string {
	pwd, _ := os.Getwd()
	return pwd
}

func GetComputerName() string {
	return os.Getenv("COMPUTERNAME")
}

func GetHostname() string {
	return os.Getenv("HOSTNAME")
}

func GetOS() string {
	return os.Getenv("OS")
}

func GetUsername() string {
	return os.Getenv("USERNAME")
}

func GetArchitecture() string {
	return os.Getenv("MSYSTEM_CARCH")
}

func GetProcessor() string {
	return os.Getenv("NUMBER_OF_PROCESSORS")
}
