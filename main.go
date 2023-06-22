package main

import (
	"flag"
	"fmt"
	"lance-light/core"
	"lance-light/render"
)

// 成功したらTrue、そうでなければFalseを返す
func writeRulesFromConfig(configFilePath string, nftableFilePath string) bool {
	rules := render.GenRulesFromConfig(configFilePath)
	core.WriteToFile(rules, nftableFilePath)
	return true
}

func exportRulesFromConfig(configFilePath string) bool {
	rules := render.GenRulesFromConfig(configFilePath)
	for _, item := range rules {
		fmt.Println(item)
	}
	return true
}

func main() {
	core.MsgInfo("LanceLight ver0.01")

	configFilePath := flag.String("f", "/etc/lance.yml", "Path of config.yml")
	nftableFilePath := flag.String("o", "/etc/nftables.lance.conf", "Path of nftables.conf")

	flag.Parse()

	operation := flag.Arg(0)

	if operation == "apply" {
		writeRulesFromConfig(*configFilePath, *nftableFilePath)
		core.ExecCommand("nft", []string{"-f", *nftableFilePath})
	} else if operation == "export" {
		exportRulesFromConfig(*configFilePath)
	}

}
