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

func ChainStart(name string) string {
	return "	chain " + name + " {"
}

func ChainEnd() string{
	return "	}"
}

func TableStart(name string) string {
	return "table inet " + name + " {"
}

func TableEnd() string {
	return "}"
}