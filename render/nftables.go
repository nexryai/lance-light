package render

import (
	"fmt"
	"lance-light/core"
	"lance-light/ip"
	"strings"
)

func MkInclude(filePath string) string {
	return fmt.Sprintf("include \"%s\"", filePath)
}

func MkDefine(name string, obj []string) string {
	joinedString := strings.Join(obj, ", ")

	rule := fmt.Sprintf("define %s = { %s };", name, joinedString)
	return rule
}

func MkBaseRules(allowed bool, direction string) string {
	policy := "drop"

	if allowed {
		policy = "accept"
	}

	return fmt.Sprintf(`		type filter hook %s priority 0; policy %s;`, direction, policy)
}

// FixMe: 関数名おかしい気がする？
func MkBaseInputRules(allowEstablished bool, allowRelated bool, allowInvalid bool) string {

	establishedRule := "drop"
	relatedRule := "drop"
	invalidRule := "drop"

	if allowEstablished {
		establishedRule = "accept"
	}

	if allowRelated {
		relatedRule = "accept"
	}

	if allowInvalid {
		invalidRule = "accept"
	}

	return fmt.Sprintf(`		ct state vmap { established : %s, related : %s, invalid : %s } `, establishedRule, relatedRule, invalidRule)
}

func MkAllowLoopbackInterface() string {
	return "\t\tiif lo accept"
}

func MkAllowPing() string {
	// ToDo: レートリミット変えられるようにするべき？
	rateLimitPerSec := 5
	return fmt.Sprintf(`		icmp type echo-request limit rate %d/second accept`, rateLimitPerSec)
}

func MkAllowIPv6Ad() string {
	// これしないとIPv6関係の接続が壊れる
	return "\t\tip6 nexthdr icmpv6 icmpv6 type { nd-neighbor-solicit, nd-router-advert, nd-neighbor-advert } accept"
}

func MkDenyIP(denyIp string) string {
	if ip.IsIPv6(denyIp) {
		return fmt.Sprintf(`		ip6 saddr %s drop`, denyIp)
	} else {
		return fmt.Sprintf(`		ip saddr %s drop`, denyIp)
	}
}

func MkAllowPort(c *core.PortConfig) string {
	rule := "\t\t"
	var allowIP string

	if c.AllowIP == "cloudflare" {
		allowIP = "$CLOUDFLARE"
	} else if c.AllowIP == "cloudflare_v6" {
		allowIP = "$CLOUDFLARE_V6"
	} else {
		allowIP = c.AllowIP
	}

	var ipVersion int8

	if allowIP == "$CLOUDFLARE_V6" || ip.IsIPv6(allowIP) {
		ipVersion = 6
	} else {
		ipVersion = 4
	}

	if c.AllowCountry != "" {
		if allowIP != "" {
			core.ExitOnError(fmt.Errorf("invalid config"), "You cannot use both allowCountry and allowIP in the same rule")
		} else {
			allowIP = fmt.Sprintf("$%s", c.AllowCountry)
		}
	}

	if c.AllowInterface != "" {
		rule += fmt.Sprintf("iifname \"%s\" ", c.AllowInterface)
	}

	rule += fmt.Sprintf("%s dport %d ", c.Proto, c.Port)

	if allowIP != "" {
		if ipVersion == 4 {
			rule += fmt.Sprintf("ip saddr %s ", allowIP)
		} else if ipVersion == 6 {
			rule += fmt.Sprintf("ip6 saddr %s ", allowIP)
		} else {
			core.GenBugCodeMessage("fea1507a-6eb7-40d4-a499-1f70ac6fd580")
		}
	}

	rule += fmt.Sprintf("accept")

	return rule
}

// ルーター用ルール
func MkAllowFwd(allowInterface string) string {
	return fmt.Sprintf("\t\tiifname %s accept", allowInterface)
}

func MkAllowForwardForNat(c *core.NatConfig) string {
	return fmt.Sprintf("\t\tiifname %s ip saddr %s accept", c.Interface, c.AllowIP)
}

func MkAllowForwardForCustomRoutes(c *core.CustomRoutesConfig) string {
	return fmt.Sprintf("\t\tiifname %s ip saddr %s ip daddr %s accept", c.AllowInterface, c.AllowIP, c.AllowDST)
}

func MkBaseRoutingRule(route string) string {
	return fmt.Sprintf("\t\ttype nat hook %s priority 100; policy accept;", route)
}

func MkMasquerade(srcIP string, outInterface string) string {
	return fmt.Sprintf("\t\tip saddr %s oifname %s masquerade", srcIP, outInterface)
}

func MkMasqueradeForNat(c *core.NatConfig) string {
	return fmt.Sprintf("\t\tiifname %s ip saddr %s masquerade", c.Interface, c.AllowIP)
}

func MkMasqueradeForCustomRoutes(c *core.CustomRoutesConfig) string {
	return fmt.Sprintf("\t\tiifname %s ip saddr %s ip daddr %s masquerade", c.AllowInterface, c.AllowIP, c.AllowDST)
}

func MkForceDNS(dnsAddress string, lanInterface string, protocol string) string {
	return fmt.Sprintf("\t\tiifname \"%s\" meta l4proto %s ip saddr != 127.0.0.1 ip daddr != %s %s dport 53 dnat to %s", lanInterface, protocol, dnsAddress, protocol, dnsAddress)
}

func MkBaseNatRule() string {
	return "\t\ttype nat hook prerouting priority dstnat;"
}

func MkNat(c *core.NatConfig) string {
	rule := fmt.Sprintf("\t\tiifname \"%s\" ", c.Interface)

	if c.AllowIP != "" {
		rule += fmt.Sprintf("ip saddr %s ", c.AllowIP)
	}

	rule += fmt.Sprintf("ip daddr { %s } %s dport %s dnat %s", c.DstIP, c.Proto, c.DstPort, c.NatTo)
	return rule
}

// ログ関係
func MkLoggingRules(policy string) string {
	return fmt.Sprintf("\t\tlog prefix \"[LanceLight] Access Denied: \" counter %s", policy)
}

// チェーンとテーブル関係
func MkChainStart(name string) string {
	return "\tchain " + name + " {"
}

func MkChainEnd() string {
	return "\t}"
}

func MkTableStart(name string) string {
	return "table inet " + name + " {"
}

func MkTableEnd() string {
	return "}"
}
