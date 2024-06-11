package iputil

import (
	"bufio"
	"fmt"
	"lance-light/internal/log"
	"net"
	"net/http"
	"regexp"
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
			log.MsgWarn("Invalid IP: " + ip)
			return false
		}
	}
	return true
}

func IsIPv6(input string) bool {
	return isValidIP(input) && strings.Contains(input, ":")
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

func FetchIpSet(url string, allowIPv6 bool) []string {
	resp, err := http.Get(url)
	log.ExitOnError(err, "failed to fetch ipset.")

	if resp.StatusCode != 200 {
		log.ExitOnError(fmt.Errorf("status code: %d", resp.StatusCode), "failed to fetch ipset.")
	}

	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	var ips []string
	for scanner.Scan() {
		i := scanner.Text()
		if !isValidIP(i) {
			log.MsgWarn("Ignore invalid ip")
		} else if !allowIPv6 && IsIPv6(i) {
			// ToDo
		} else {
			ips = append(ips, i)
		}
	}

	return ips
}

func ExtractIPAddress(input string) (string, error) {
	// IPアドレスの正規表現パターン
	ipPattern := `(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`

	// 正規表現にマッチする部分を抽出
	re := regexp.MustCompile(ipPattern)
	matches := re.FindStringSubmatch(input)

	// マッチが見つからなかった場合
	if len(matches) < 2 {
		return "", fmt.Errorf("no ip found")
	}

	if !isValidIP(matches[1]) {
		return "", fmt.Errorf("invalid ip found")
	}

	// 抽出したIPアドレスを返す
	return matches[1], nil
}
