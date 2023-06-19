package render

import (
	"fmt"
)

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

/*
func SecureRules(denyTorIPs bool, denyAbuseIPs bool, denyPublicProxyIPs bool, alwaysDenyIPs []string{}, alwaysDenyASNs []int{}) {
	//#ToDo
}

func PortRules(port int, allowIPs string, allowInterface string, allowProto string) string {
	//#ToDo
}


func RouterRules(lanInterface string, wanInterface string, forceDNS string) {
	//#ToDo
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
