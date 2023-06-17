package render

import (
	"as4guard/core"
	"io/ioutil"
	"net/http"
	"strings"
	"strconv"
)

func getCloudflareIPs(version int) []string {

	if version != 4 && version != 6 {
		core.MsgErr("Internal error. EUID:26987ba0-2355-418b-9bc8-c0d76189cd16 \nPlease contact the developer.")
	}

	resp, err := http.Get("https://www.cloudflare.com/ips-v" + strconv.Itoa(version))
	if err != nil {
		core.MsgErr("Error fetching IP list:")
		return []string{}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		core.MsgErr("Error reading response body:")
		return []string{}
	}

	// レスポンスボディを文字列に変換し、改行文字で分割してリストに代入
	ipList := strings.Split(string(body), "\n")
	return ipList
}


func baseRules(allowOut bool, allowIn bool, allowFwd bool) string {
	#ToDo
}

func secureRules(denyTorIPs bool, denyAbuseIPs bool, denyPublicProxyIPs bool, alwaysDenyIPs []string{}, alwaysDenyASNs []int{}) {
	#ToDo
}

func portRules(port int, allowIPs string, allowInterface string, allowProto string) string {
	#ToDo
}


func routerRules(lanInterface string, wanInterface string, forceDNS string) {
	#ToDo
}
