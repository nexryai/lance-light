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

func MkAllowPing() string {
	// ToDo: レートリミット変えられるようにするべき？
	rateLimitPerSec := 5
	return fmt.Sprintf(`		icmp type echo-request limit rate %d second accept`, rateLimitPerSec)
}

func MkAllowIPv6Ad() string {
	return "		ip6 nexthdr icmpv6 icmpv6 type { nd-neighbor-solicit, nd-router-advert, nd-neighbor-advert } accept"
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

func MkChainEnd() string{
	return "	}"
}

func MkTableStart(name string) string {
	return "table inet " + name + " {"
}

func MkTableEnd() string {
	return "}"
}