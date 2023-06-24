package render

import (
	"fmt"
	"strings"
)

func MkFlushTable(tableName string) string {
	return fmt.Sprintf("flush table inet %s", tableName)
}

func MkDefine(name string, obj []string) string {
	joinedString := strings.Join(obj, ", ")

	rule := fmt.Sprintf("define %s = { %s }", name, joinedString)
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
	return fmt.Sprintf(`		ip saddr %s drop`, denyIp)
}

func MkAllowPort(port int, allowIP string, allowInterface string, allowProto string) string {
	rule := "\t\t"

	if allowInterface != "" {
		rule += fmt.Sprintf("iifname \"%s\" ", allowInterface)
	}

	rule += fmt.Sprintf("%s dport %d ", allowProto, port)

	if allowIP != "" {
		rule += fmt.Sprintf("ip saddr %s ", allowIP)
	}

	rule += fmt.Sprintf("accept")

	return rule
}

// ルーター用ルール
func MkAllowFwd(allowInterface string) string {
	return fmt.Sprintf("\t\tiifname %s accept", allowInterface)
}

func MkBaseRoutingRule(route string) string {
	return fmt.Sprintf("\t\ttype nat hook %s priority 100; policy accept;", route)
}

func MkMasquerade(srcIP string, outInterface string) string {
	return fmt.Sprintf("\t\tip saddr %s oifname %s masquerade", srcIP, outInterface)
}

func MkForceDNS(dnsAddress string, lanInterface string, protocol string) string {
	return fmt.Sprintf("\t\tiifname \"%s\" meta l4proto %s ip saddr != 127.0.0.1 ip daddr != %s %s dport 53 dnat to %s", lanInterface, protocol, dnsAddress, protocol, dnsAddress)
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
