package render

import (
	"errors"
	"github.com/lorenzosaino/go-sysctl"
	"io/ioutil"
	"lance-light/core"
	"lance-light/ip"
	"net/http"
	"os"
	"strconv"
	"strings"
)

/*
nftablesルールをレンダリングする。基本的に1行の内容を1つづつ配列に格納して返す
*/

func shouldGenPreroutingRules(config *core.Config) bool {
	if config.Router.ConfigAsRouter && config.Router.ForceDNS != "" {
		// ForceDNSが設定されているならtrue
		return true
	} else if len(config.Nat) != 0 {
		// Nat設定があるならtrue
		return true
	} else {
		return false
	}
}

func getCloudflareIPs(version int) ([]string, error) {

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
	if !ip.CheckIPAddresses(cfIpList) {
		core.ExitOnError(errors.New("invalid IP from API"), core.GenBugCodeMessage("8a04693b-9a36-422b-81b6-2270ad8e357b"))
	}

	return cfIpList, nil
}

func GenIpDefineRules(rule string, config *core.Config) ([]string, error) {
	rules := []string{}

	if rule == "cloudflare" {
		// CloudflareのIPを取得し定義する。
		var clouflareIPsV4 []string
		var clouflareIPsV6 []string
		var e error

		clouflareIPsV4, e = getCloudflareIPs(4)
		if e != nil {
			return rules, e
		}

		if config.Default.EnableIPv6 {
			clouflareIPsV6, e = getCloudflareIPs(6)
			if e != nil {
				return rules, e
			}
		}

		rules = append(rules, MkDefine("CLOUDFLARE", clouflareIPsV4), MkDefine("CLOUDFLARE_V6", clouflareIPsV6))

	} else {
		core.ExitOnError(errors.New("unexpected argument"), core.GenBugCodeMessage("16ee8518-2ad6-4946-8d10-fbc77a1da586"))
	}

	return rules, nil
}

func GenRulesFromConfig(config *core.Config) []string {

	rules := []string{}

	// IpDefineFilePathをincludeする
	// IpDefineFilePathにはCloudflareのIPやAubseIPがキャッシュされている
	rules = append(rules, MkInclude(config.Nftables.IpDefineFilePath))

	//テーブル作成
	rules = append(rules, MkTableStart("lance"))

	if config.Default.AllowAllIn {
		core.MsgWarn("Input is allowed by default. This is a VERY UNSAFE setting. You MUST not use this setting unless you know what you are doing.")
	}

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
	for _, r := range config.Ports {
		if r.Proto == "" {
			// 本当はprotoを必須にしたいけど互換性維持のため
			r.Proto = "tcp"
			rules = append(rules, MkAllowPort(&r))
			r.Proto = "udp"
			rules = append(rules, MkAllowPort(&r))
		} else {
			rules = append(rules, MkAllowPort(&r))
		}
	}

	// ログが有効ならログする
	if config.Default.EnableLogging {
		rules = append(rules, MkLoggingRules("drop"))
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

	// ルーターとして構成するならLAN→WANへのマスカレードを許可する
	if config.Router.ConfigAsRouter {
		rules = append(rules, MkBaseInputRules(true, true, false))

		for _, lanInterface := range config.Router.LANInterfaces {
			rules = append(rules, MkAllowFwd(lanInterface))
		}
	}

	// ポート転送構成時にForwardを許可する
	if len(config.Nat) != 0 {
		for _, r := range config.Nat {
			rules = append(rules, MkAllowForwardForNat(&r))
		}
	}

	rules = append(rules, MkChainEnd())

	// POSTROUTINGチェーン
	if config.Router.ConfigAsRouter || len(config.Nat) != 0 {
		sysctlIpForward, err := sysctl.Get("net.ipv4.ip_forward")

		if err != nil {
			core.MsgWarn("Failed to get sysctl value")
		} else if sysctlIpForward == "0" {
			core.MsgWarn("net.ipv4.ip_forward is set to 0.")
		}

		rules = append(rules,
			MkChainStart("postrouting"),
			MkBaseRoutingRule("postrouting"))

		// ルーターとして構成するときのLAN→WANのマスカレード
		if config.Router.ConfigAsRouter {
			for _, privateNetworkAddress := range config.Router.PrivateNetworkAddresses {
				rules = append(rules, MkMasquerade(privateNetworkAddress, config.Router.WANInterface))
			}
		}

		// ポート転送有効時のマスカレード
		if len(config.Nat) != 0 {
			for _, r := range config.Nat {
				rules = append(rules, MkMasqueradeForNat(&r))
			}
		}

		rules = append(rules, MkChainEnd())
	}

	// PREROUTINGチェーン
	if shouldGenPreroutingRules(config) {
		rules = append(rules, MkChainStart("prerouting"))

		if config.Router.ConfigAsRouter {
			rules = append(rules, MkBaseRoutingRule("prerouting"))
		} else if len(config.Nat) != 0 {
			rules = append(rules, MkBaseNatRule())
		}

		if config.Router.ForceDNS != "" {
			for _, lanInterface := range config.Router.LANInterfaces {
				rules = append(rules, MkForceDNS(config.Router.ForceDNS, lanInterface, "udp"))
				rules = append(rules, MkForceDNS(config.Router.ForceDNS, lanInterface, "tcp"))
			}
		}

		// ポート転送有効時のNAT構成
		if len(config.Nat) != 0 {
			for _, r := range config.Nat {
				rules = append(rules, MkNat(&r))
			}
		}

		rules = append(rules, MkChainEnd())
	}

	// テーブル終了
	rules = append(rules, MkTableEnd())

	return rules
}
