package ip

import (
	"bufio"
	"lance-light/core"
	"net"
	"net/http"
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

// stringが有効なIPアドレスか識別する
func isValidIP(ip string) bool {
	if _, _, err := net.ParseCIDR(ip); err == nil {
		return true
	}
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}

// []string内のすべての文字が有効なIPアドレスか識別する
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

// AbuseIPDBに通報可能なIPかどうか
func IsReportableAddress(ip string) bool {
	if !isValidIP(ip) {
		return false
	} else if isPrivateAddress(ip) {
		return false
	} else {
		return true
	}
}

func FetchIpSet(url string) []string {
	resp, err := http.Get(url)
	core.ExitOnError(err, "failed to fetch ipset.")

	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	var ips []string
	for scanner.Scan() {
		i := scanner.Text()
		if isValidIP(i) && !IsIPv6(i) {
			ips = append(ips, i)
		} else {
			//core.MsgWarn("FetchIpSet: Ignore invalid line.")
		}
	}

	return ips
}
