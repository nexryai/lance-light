package render

import (
	"fmt"
	"github.com/lorenzosaino/go-sysctl"
	"lance-light/core"
	"lance-light/entities"
	"lance-light/ip"
)

/*
nftablesルールをレンダリングする。基本的に1行の内容を1つづつ配列に格納して返す
*/

func containsString(arr []string, target string) bool {
	for _, element := range arr {
		if element == target {
			return true
		}
	}
	return false
}

func shouldDefineCloudflareIPs(config *core.Config) bool {
	for _, p := range config.Ports {
		if p.AllowIP == "cloudflare" || p.AllowIP == "cloudflare_v6" {
			return true
		}
	}

	return false
}

func GenIpDefineRules(config *core.Config) ([]string, error) {
	rules := []string{}

	// CloudflareのIPを取得し定義する
	if shouldDefineCloudflareIPs(config) {
		var clouflareIPsV4 []string
		var clouflareIPsV6 []string

		clouflareIPsV4 = ip.FetchIpSet("https://www.cloudflare.com/ips-v4", false)

		if config.Default.EnableIPv6 {
			clouflareIPsV6 = ip.FetchIpSet("https://www.cloudflare.com/ips-v6", true)
		}

		rules = append(rules, MkDefine("CLOUDFLARE", clouflareIPsV4), MkDefine("CLOUDFLARE_V6", clouflareIPsV6))
	}

	// AllowCountryに存在する国コードのIPを取得し定義する
	var countries []string
	seen := make(map[string]bool)

	for _, p := range config.Ports {
		c := p.AllowCountry
		if !seen[c] && c != "" {
			countries = append(countries, c)
			seen[c] = true
		}
	}

	for _, c := range countries {
		url := fmt.Sprintf("https://www.ipdeny.com/ipblocks/data/countries/%s.zone", c)
		r := MkDefine(c, ip.FetchIpSet(url, false))
		rules = append(rules, r)
	}

	// ユーザー定義のipsetをロードする
	for _, s := range config.IpSet {
		var addrs []string

		if s.Url != "" {
			addrs = append(addrs, ip.FetchIpSet(s.Url, false)...)
		}

		addrs = append(addrs, s.Ip...)
		rules = append(rules, MkDefine(s.Name, addrs))
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

	// AlwaysDenyTorならTorのIPを拒否
	// Torの出口のIPは頻繁に変動するため、将来的にキャッシュの対象となるGenIpDefineRulesで生成せずここで毎回取得する
	if config.Security.AlwaysDenyTor {
		for _, denyIP := range ip.FetchIpSet("https://check.torproject.org/torbulkexitlist?ip=1.1.1.1", false) {
			alwaysDenyIP = append(alwaysDenyIP, denyIP)
		}
	}

	// alwaysDenyIPを拒否
	for _, denyIP := range alwaysDenyIP {
		rules = append(rules, MkDenyIP(denyIP))
	}

	// pingを許可するなら許可
	if config.Default.AllowPing {
		rules = append(rules, MkAllowPing(), MkAllowPingICMPv6())
	}

	if config.Security.AlwaysDenyAbuseIP {
		core.MsgDebug("Always Deny AbuseIP")
	}

	// IPv6関係
	if config.Default.EnableIPv6 {
		rules = append(rules, MkAllowIPv6Ad())
	}

	// SYN-floodレートリミット
	rules = append(rules, MkJumpToSynFloodLimiter())

	// 許可したポートをallow
	for _, r := range config.Ports {
		if r.AllowIP != "" && r.AllowCountry != "" {
			// AllowCountryが設定されているとAllowIPが上書きされてしまうので対策
			r.AllowIP = fmt.Sprintf("{ $%s, %s }", r.AllowCountry, r.AllowIP)
			r.AllowCountry = ""
		}

		if r.Proto == "" {
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
		MkBaseRules(config.Default.AllowAllOut, "output"))

	if !config.Default.AllowAllOut {
		// ICMPとループバックは許可
		rules = append(rules, MkAllowLoopbackInterface(),
			MkAllowIcmpOutgoing(),
			MkAllowLocalhostOutgoing())

		// IPv6が有効ならIPv6のICMPも許可
		if config.Default.EnableIPv6 {
			rules = append(rules, MkAllowIcmpv6Outgoing())
		}

		for _, r := range config.Outgoing.Allowed {
			if r.Proto == "" {
				r.Proto = "tcp"
				rules = append(rules, MkAllowOutgoing(&r))
				r.Proto = "udp"
				rules = append(rules, MkAllowOutgoing(&r))
			} else {
				rules = append(rules, MkAllowOutgoing(&r))
			}
		}

		// Tailscaleと併用できるようにする (https://tailscale.com/kb/1082/firewall-ports)
		//  - 41641への発信と3478への発信を許可する
		if containsString(config.Outgoing.Compatibility, "tailscale") {
			rules = append(rules, MkAllowOutgoing(&core.OutgoingAllowConfig{
				Dport: "41641",
				Proto: "udp",
			}), MkAllowOutgoing(&core.OutgoingAllowConfig{
				Dport: "3478",
				Proto: "udp",
			}))
		}

		// Cloudflaredと併用できるようにする
		// https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/deploy-tunnels/tunnel-with-firewall/
		if containsString(config.Outgoing.Compatibility, "cloudflare_tunnel") {
			rules = append(rules, MkAllowOutgoing(&core.OutgoingAllowConfig{
				Dport: "7844",
				Proto: "tcp",
			}), MkAllowOutgoing(&core.OutgoingAllowConfig{
				Dport: "7844",
				Proto: "udp",
			}), MkAllowOutgoing(&core.OutgoingAllowConfig{
				Dport: "443",
				Proto: "tcp",
			}))
		}

		// SYNじゃない（=こっちからの発信じゃない）なら許可
		rules = append(rules, MkAllowNonSynOutgoing())

		// ログが有効ならログする
		if config.Default.EnableLogging {
			rules = append(rules, MkLoggingForOutgoing())
		}
	}

	rules = append(rules, MkChainEnd())

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

		// カスタムルート設定時にForward許可する
		for _, r := range config.Router.CustomRoutes {
			rules = append(rules, MkAllowForwardForCustomRoutes(&r))
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

		// DNATするときの戻り通信用SNAT
		if len(config.Nat) != 0 {
			for _, c := range config.Nat {
				internalIP, err := ip.ExtractIPAddress(c.NatTo)
				if err != nil {
					panic("invalid ip in config")
				}

				snat := entities.SnatForDnat{
					ExternalInterface: c.Interface,
					InternalIP:        internalIP,
					ExternalIP:        c.DstIP,
				}
				rules = append(rules, MkSnatForDnat(&snat))
			}
		}

		if config.Router.ConfigAsRouter {
			// ルーターとして構成するときのLAN→WANのマスカレード
			for _, privateNetworkAddress := range config.Router.PrivateNetworkAddresses {
				rules = append(rules, MkMasquerade(privateNetworkAddress, config.Router.WANInterface))
			}

			// カスタムルート設定時のマスカレード設定
			for _, r := range config.Router.CustomRoutes {
				rules = append(rules, MkMasqueradeForCustomRoutes(&r))
			}
		}

		rules = append(rules, MkChainEnd())
	}

	// PREROUTINGチェーン
	rules = append(rules, MkChainStart("prerouting"))

	if config.Router.ConfigAsRouter {
		rules = append(rules, MkBaseRoutingRule("prerouting"))
	} else if len(config.Nat) != 0 {
		rules = append(rules, MkBaseNatRule())
	}

	// 不正なパケットととりあえず全部弾くべき攻撃を遮断
	// inputチェーンよりpreroutingの方が優先されるのでここに入れる
	rules = append(rules, MkDropInvalid())
	if !config.Security.DisablePortScanProtection {
		rules = append(rules, MkBlockTcpXmas(), MkBlockTcpNull(), MkBlockTcpMss())
	} else {
		core.MsgWarn("Port scan protection is DISABLED!")
	}

	if !config.Security.DisableIpFragmentsBlock {
		rules = append(rules, MkBlockIPFragments())
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

	// SYN-flood対策
	rules = append(rules, MkChainStart("syn-flood"),
		MkRateLimit(30, 60, "SYN-flood"),
		MkChainEnd())

	// テーブル終了
	rules = append(rules, MkTableEnd())

	return rules
}
