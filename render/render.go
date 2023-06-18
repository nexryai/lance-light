package render

import (
	"lance-light/core"
	"lance-light/ip"
	"io/ioutil"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
)

/*
nftableルールをレンダリングする。基本的に1行の内容を1つづつ配列に格納して返す
*/

func getCloudflareIPs(version int) []string {

	if version != 4 && version != 6 {
		core.MsgErr("Internal error. EUID:26987ba0-2355-418b-9bc8-c0d76189cd16 \nPlease contact the developer.")
		os.Exit(2)
	}

	resp, err := http.Get("https://www.cloudflare.com/ips-v" + strconv.Itoa(version))
	core.ExitOnError(err, "Failed to fetch Cloudflare's list of IP addresses. If checking your network connection does not resolve the issue, please contact the developer.")
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	core.ExitOnError(err, "An unexpected error occurred while retrieving Cloudflare's IP address. The request was successful, but an error occurred while reading the response body.")

	// レスポンスボディを文字列に変換し、改行文字で分割してリストに代入
	CfIpList := strings.Split(string(body), "\n")
	// CfIpList=[]string{"192.168.0.1", "10.0.0.1", "256.0.0.1", "172.16.0.1"}

	if !ip.CheckIPAddresses(CfIpList) {
		core.ExitOnError(errors.New("Invalid IP"), "An error occurred while retrieving the IP list from Cloudflare. The request was successful, but an invalid IP address was detected.")
	}

	return CfIpList
}


func GenRulesFromConfig(configFilePath string) []string {
	core.LoadConfig(configFilePath)
	return getCloudflareIPs(4)
}