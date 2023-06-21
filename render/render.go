package render

import (
	"errors"
	"io/ioutil"
	"lance-light/core"
	"lance-light/ip"
	"net/http"
	"os"
	"strconv"
	"strings"
)

/*
nftableルールをレンダリングする。基本的に1行の内容を1つづつ配列に格納して返す
*/

func getCloudflareIPs(version int) []string {

	if version != 4 && version != 6 {
		core.MsgErr("Internal error. EUID:26987ba0-2355-418b-9bc8-c0d76189cd16 \nPlease contact the developer.")
		os.Exit(2)
	}

	resp, err := http.Get("https://www.cloudflare.com/ips-v" + strconv.Itoa(version))
	core.ExitOnError(err, "Failed to fetch Cloudflare's list of IP addresses. If checking your network connection does not resolve the issue, please contact the developer.")
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	core.ExitOnError(err, "An unexpected error occurred while retrieving Cloudflare's IP address. The request was successful, but an error occurred while reading the response body.")

	// レスポンスボディを文字列に変換し、改行文字で分割してリストに代入
	cfIpList := strings.Split(string(body), "\n")
	//CfIpList=[]string{"192.168.0.1", "10.0.0.1", "256.0.0.1", "172.16.0.1"}

	if !ip.CheckIPAddresses(cfIpList) {
		core.ExitOnError(errors.New("invalid IP from API"), "An error occurred while retrieving the IP list from Cloudflare. The request was successful, but an invalid IP address was detected.")
	}

	return cfIpList
}

func getAllCloudflareIPs() []string {
	cfAllIpList := getCloudflareIPs(4)
	cfAllIpList = append(cfAllIpList, getCloudflareIPs(6)...)
	return cfAllIpList
}

func GenRulesFromConfig(configFilePath string) []string {
	config := core.LoadConfig(configFilePath)

	rules := []string{}

	// CloudflareのIPを取得し定義する。
	if config.Default.EnableIPv6 {
		rules = append(rules, MkDefine("CLOUDFLARE", getAllCloudflareIPs()))
	} else {
		rules = append(rules, MkDefine("CLOUDFLARE", getCloudflareIPs(4)))
	}

	//テーブル作成
	rules = append(rules, MkTableStart("filter"))

	// INPUTルール作成
	rules = append(rules,
		MkChainStart("input"),
		MkBaseRules(config.Default.AllowAllIn, "input"))

	// これは変えられるようにするべき？
	rules = append(rules,
		MkBaseInputRules(true, true, false),
		MkAllowLoopbackInterface())

	alwaysDenyIP := []string{}

	alwaysDenyIP = append(alwaysDenyIP, config.Security.AlwaysDenyIP...)

	// alwaysDenyASNをIPのCIDRに変換
	for _, denyASN := range config.Security.AlwaysDenyASN {
		alwaysDenyIP = append(alwaysDenyIP, ip.GetIpRangeFromASN(denyASN)...)
	}

	// alwaysDenyIPを拒否
	for _, denyIP := range alwaysDenyIP {
		rules = append(rules, MkDenyIP(denyIP))
	}

	// pingを許可するなら許可
	if config.Default.AllowPing {
		rules = append(rules, MkAllowPing())
	}

	if config.Security.AlwaysDenyAbuseIP {
		core.MsgDebug("Always Deny AbuseIP")
	}

	// IPv6関係
	if config.Default.EnableIPv6 {
		rules = append(rules, MkAllowIPv6Ad())
	}

	// 許可したポートをallow
	for _, allowPort := range config.Ports {
		var allowIP string

		if allowPort.AllowIP == "cloudflare" {
			allowIP = "$CLOUDFLARE"
		} else {
			allowIP = allowPort.AllowIP
		}

		if allowPort.Proto == "" {
			rules = append(rules, MkAllowPort(allowPort.Port, allowIP, allowPort.AllowInterface, "tcp"))
			rules = append(rules, MkAllowPort(allowPort.Port, allowIP, allowPort.AllowInterface, "udp"))
		} else {
			rules = append(rules, MkAllowPort(allowPort.Port, allowIP, allowPort.AllowInterface, allowPort.Proto))
		}
	}

	// INPUTチェーン終了
	rules = append(rules, MkChainEnd())

	// OUTPUTチェーン
	rules = append(rules, MkChainStart("output"),
		MkBaseRules(config.Default.AllowAllOut, "output"),
		MkChainEnd())

	// FORWARDチェーン
	rules = append(rules, MkChainStart("forward"))

	if config.Default.AllowAllFwd {
		core.MsgWarn("Forwarding is allowed by default. This is an unsafe setting and you usually don't need to do this.")
	}

	rules = append(rules, MkBaseRules(config.Default.AllowAllFwd, "forward"))

	if config.Router.ConfigAsRouter {
		rules = append(rules,
			MkBaseInputRules(true, true, false),
			MkAllowFwd(config.Router.LANInterface))
	}

	rules = append(rules, MkChainEnd())

	// POSTROUTINGチェーン
	if config.Router.ConfigAsRouter {
		rules = append(rules,
			MkChainStart("postrouting"),
			MkBasePostroutingRule(),
			MkMasquerade(config.Router.PrivateNetworkAddress, config.Router.WANInterface),
			MkChainEnd())
	}

	// テーブル終了
	rules = append(rules, MkTableEnd())

	return rules
}
