package render

import (
	"fmt"
	"strings"
)

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

	return fmt.Sprintf(`		type filter hook %s priority 0; %s;`, direction, policy)
}

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
	return "		iif lo accept"
}

func MkAllowPing() string {
	// ToDo: レートリミット変えられるようにするべき？
	rateLimitPerSec := 5
	return fmt.Sprintf(`		icmp type echo-request limit rate %d second accept`, rateLimitPerSec)
}

func MkAllowIPv6Ad() string {
	return "		ip6 nexthdr icmpv6 icmpv6 type { nd-neighbor-solicit, nd-router-advert, nd-neighbor-advert } accept"
}

func MkDenyIP(denyIp string) string {
	return fmt.Sprintf(`		ip saddr %s drop`, denyIp)
}

func MkAllowPort(port int, allowIP string, allowInterface string, allowProto string) string {
	rule := "		"

	if allowInterface != "" {
		rule += fmt.Sprintf("iifname \"%s\" ", allowInterface)
	}

	rule += fmt.Sprintf("%s ", allowProto)

	if allowIP != "" {
		rule += fmt.Sprintf("saddr %s ", allowIP)
	}

	rule += fmt.Sprintf("dport %d accept", port)

	return rule
}

// ルーター用ルール
func MkAllowFwd(allowInterface string) string {
	return fmt.Sprintf("\t\tiifname %s accept", allowInterface)
}

func MkBasePostroutingRule() string {
	return "\t\ttype nat hook postrouting priority 100; policy accept;"
}

func MkMasquerade(srcIP string, outInterface string) string {
	return fmt.Sprintf("\t\tip saddr %s oifname %s masquerade", srcIP, outInterface)
}

/*
func MkForceDNS(dnsAddress string) string {
	//ToDo
}
*/

func MkChainStart(name string) string {
	return "	chain " + name + " {"
}

func MkChainEnd() string {
	return "	}"
}

func MkTableStart(name string) string {
	return "table inet " + name + " {"
}

func MkTableEnd() string {
	return "}"
}
