package render

import (
	"lance-light/core"
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


func GenRulesFromConfig()