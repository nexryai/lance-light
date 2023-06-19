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
	CfIpList := strings.Split(string(body), "\n")
	//CfIpList=[]string{"192.168.0.1", "10.0.0.1", "256.0.0.1", "172.16.0.1"}

	if !ip.CheckIPAddresses(CfIpList) {
		core.ExitOnError(errors.New("invalid IP from API"), "An error occurred while retrieving the IP list from Cloudflare. The request was successful, but an invalid IP address was detected.")
	}

	return CfIpList
}

func GenRulesFromConfig(configFilePath string) []string {
	config := core.LoadConfig(configFilePath)

	rules := []string{}

	//テーブル作成
	rules = append(rules, MkTableStart("filter"))

	// INPUTルール作成（ToDo: ポート許可）
	rules = append(rules, MkChainStart("input"))
	rules = append(rules, MkBaseRules(config.Default.AllowAllIn, "input"))

	// これは変えられるようにするべき？
	rules = append(rules, MkBaseInputRules(true, true, false))

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

	// IPv6関係（ToDo: IPv6が無効なら追加しない）
	rules = append(rules, MkAllowIPv6Ad())

	// INPUTチェーン終了
	rules = append(rules, MkChainEnd())

	// テーブル終了
	rules = append(rules, MkTableEnd())

	return rules
}
