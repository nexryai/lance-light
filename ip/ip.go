package ip

import (
	"lance-light/core"
	"net"
	"strings"
)

func isPrivateAddress(address string) bool {
	// 0.0.0.0は127.0.0.1と扱うOSが多いのでプライベートアドレスとして扱う
	if address == "0.0.0.0" {
		return true
	}

	ip := net.ParseIP(address)
	return ip != nil && (ip.IsLoopback() || ip.IsPrivate())
}

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

func IsReportableAddress(ip string) bool {
	if !isValidIP(ip) {
		return false
	} else if isPrivateAddress(ip) {
		return false
	} else {
		return true
	}
}
