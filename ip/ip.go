package ip

import (
	"lance-light/core"
	"net"
	"strings"
)

func isValidIP(ip string) bool {
	if _, _, err := net.ParseCIDR(ip); err == nil {
		return true
	}
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}

func CheckIPAddresses(ipAddresses []string) bool {
	for _, ip := range ipAddresses {
		if !isValidIP(ip) {
			core.MsgWarn("Invalid IP: " + ip)
			return false
		}
	}
	return true
}

func IsIPv6(input string) bool {
	return strings.Contains(input, ":")
}
