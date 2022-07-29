package gipv4

import (
	"net"
	"strings"
)

// GetHostByName returns the IPv4 address corresponding to a given Internet host name.
func GetHostByName(hostname string) (string, error) {
	ips, err := net.LookupIP(hostname)
	if ips != nil {
		for _, v := range ips {
			if v.To4() != nil {
				return v.String(), nil
			}
		}
		return "", nil
	}
	return "", err
}

// GetHostsByName returns a list of IPv4 addresses corresponding to a given Internet
// host name.
func GetHostsByName(hostname string) ([]string, error) {
	ips, err := net.LookupIP(hostname)
	if ips != nil {
		var ipStrings []string
		for _, v := range ips {
			if v.To4() != nil {
				ipStrings = append(ipStrings, v.String())
			}
		}
		return ipStrings, nil
	}
	return nil, err
}

// GetNameByAddr returns the Internet host name corresponding to a given IP address.
func GetNameByAddr(ipAddress string) (string, error) {
	names, err := net.LookupAddr(ipAddress)
	if names != nil {
		return strings.TrimRight(names[0], "."), nil
	}
	return "", err
}
