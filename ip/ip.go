package ip

import (
	"errors"
	"io/ioutil"
	"lance-light/core"
	"net"
	"net/http"
	"os"
	"strconv"
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

// CloudflareのサーバーのIPリストを取得する
func GetCloudflareIPs(version int) ([]string, error) {

	if version != 4 && version != 6 {
		core.MsgErr("Internal error. EUID:26987ba0-2355-418b-9bc8-c0d76189cd16 \nPlease contact the developer.")
		os.Exit(2)
	}

	resp, err := http.Get("https://www.cloudflare.com/ips-v" + strconv.Itoa(version))
	defer resp.Body.Close()

	// ネットワークエラーならExitOnErrorしない
	if err != nil {
		core.MsgErr("Failed to fetch Cloudflare's list of IP addresses. If checking your network connection does not resolve the issue, please contact the developer.")
		return []string{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	core.ExitOnError(err, "An unexpected error occurred while retrieving Cloudflare's IP address. The request was successful, but an error occurred while reading the response body.")

	// レスポンスボディを文字列に変換し、改行文字で分割してリストに代入
	cfIpList := strings.Split(string(body), "\n")

	// 取得したIPが正しいか念の為確認する
	if !CheckIPAddresses(cfIpList) {
		core.ExitOnError(errors.New("invalid IP from API"), core.GenBugCodeMessage("8a04693b-9a36-422b-81b6-2270ad8e357b"))
	}

	return cfIpList, nil
}
