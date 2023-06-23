package main

import (
	"flag"
	"fmt"
	"lance-light/core"
	"lance-light/render"
)

// 成功したらTrue、そうでなければFalseを返す
func writeRulesFromConfig(configFilePath string, nftableFilePath string, addFlushRule bool) bool {
	rules := render.GenRulesFromConfig(configFilePath, addFlushRule)
	core.WriteToFile(rules, nftableFilePath)
	return true
}

func exportRulesFromConfig(configFilePath string) bool {
	rules := render.GenRulesFromConfig(configFilePath, false)
	for _, item := range rules {
		fmt.Println(item)
	}
	return true
}

func main() {
	configFilePath := flag.String("f", "/etc/lance.yml", "Path of config.yml")
	nftableFilePath := flag.String("o", "/etc/nftables.lance.conf", "Path of nftables.conf")

	flag.Parse()

	operation := flag.Arg(0)

	if operation == "apply" {

		// 設定をリセットして再設定する
		writeRulesFromConfig(*configFilePath, *nftableFilePath, true)
		core.ExecCommand("nft", []string{"-f", *nftableFilePath})
		core.MsgInfo("Firewall settings have been applied successfully.")

	} else if operation == "enable" {

		// 設定を適用する
		writeRulesFromConfig(*configFilePath, *nftableFilePath, false)
		core.ExecCommand("nft", []string{"-f", *nftableFilePath})
		core.MsgInfo("LanceLight firewall is enabled.")

	} else if operation == "export" {

		// エクスポート
		exportRulesFromConfig(*configFilePath)

	} else if operation == "disable" {

		// 設定をアンロードする
		core.ExecCommand("nft", []string{"flush", "table", "inet", "lance"})
		core.MsgInfo("LanceLight firewall is disabled.")

	}

}
