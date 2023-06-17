package render

/*
func BaseRules(allowOut bool, allowIn bool, allowFwd bool) string {
	//#ToDo
}

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

func ChainStart(name) {
	return "	chain " + name + " {"
}

func ChainEnd() {
	return "	}"
}

func TableStart(name) {
	return "table inet " + name + " {"
}

func TableEnd() {
	return "}"
}