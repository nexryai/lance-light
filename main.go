package main

import (
	"flag"
	"fmt"
	"lance-light/core"
	"lance-light/render"
	"os"
)

func applyNftablesRules(configFilePath string) {
	core.ExecCommand("nft", []string{"-f", configFilePath})
}

func writeRulesFromConfig(config *core.Config, addFlushRule bool) bool {

	ipDefineRules, err := render.GenIpDefineRules("cloudflare", config)
	if err != nil {
		//ToDo: ipDefineFilePathが存在しなければ失敗扱い
	} else {
		core.WriteToFile(ipDefineRules, config.Nftables.IpDefineFilePath)
	}

	rules := render.GenRulesFromConfig(config, addFlushRule)
	core.WriteToFile(rules, config.Nftables.NftablesFilePath)
	return true
}

func exportRulesFromConfig(config *core.Config) bool {

	rules := render.GenRulesFromConfig(config, false)
	for _, item := range rules {
		fmt.Println(item)
	}
	return true
}

func showHelp() {
	fmt.Println("LanceLight firewall - Yet another human-friendly firewall \n\n",
		"(c)2023 nexryai\nThis program is licensed under the Mozilla Public License Version 2.0, and anyone can audit and contribute to it.\n\n\n",
		"[usage]\n",
		"Enable firewall:\n  ▶ llfctl enable\n\n",
		"Apply rules when configuration is updated:\n  ▶ llfctl apply\n\n",
		"Disable firewall:\n  ▶ llfctl disable\n\n",
		"[options]\n",
		"-f [PATH]: Specify the path to the configuration file (Default: /etc/lance.yml)\n")
}

func main() {
	configFilePath := flag.String("f", "/etc/lance.yml", "Path of config.yml")
	debugMode := flag.Bool("debug", false, "Enable debug mode")

	flag.Parse()

	// 現状デバッグモードはログ以外に影響を与えないが将来的に変わる可能性もある。
	// 環境変数以外にもいいやり方あるかもしれない。毎回configをcore.MsgDbg呼ぶ時に引数渡すのは面倒だから避けたい。
	if *debugMode {
		core.MsgInfo("debug mode!")
		os.Setenv("LANCE_DEBUG_MODE", "true")
	} else {
		os.Setenv("LANCE_DEBUG_MODE", "false")
	}

	core.MsgDebug("configFilePath: " + *configFilePath)

	config := core.LoadConfig(*configFilePath)

	if config.Nftables.NftablesFilePath == "" {
		config.Nftables.NftablesFilePath = "/etc/nftables.lance.conf"
	}

	if config.Nftables.IpDefineFilePath == "" {
		config.Nftables.IpDefineFilePath = "/etc/nftables.ipdefine.conf"
	}

	operation := flag.Arg(0)

	if operation == "apply" {

		// 設定をリセットして再設定する
		writeRulesFromConfig(&config, true)

		// nftコマンドを実行して適用
		applyNftablesRules(config.Nftables.NftablesFilePath)

		core.MsgInfo("Firewall settings have been applied successfully.")

	} else if operation == "enable" {

		// 設定を適用する
		writeRulesFromConfig(&config, false)

		// nftコマンドを実行して適用
		applyNftablesRules(config.Nftables.NftablesFilePath)

		core.MsgInfo("LanceLight firewall is enabled.")

	} else if operation == "export" {

		// エクスポート
		core.MsgDebug(fmt.Sprintf("configFilePath: %s", *configFilePath))
		exportRulesFromConfig(&config)

	} else if operation == "disable" {

		// 設定をアンロードする
		core.ExecCommand("nft", []string{"flush", "table", "inet", "lance"})
		core.MsgInfo("LanceLight firewall is disabled.")

	} else if operation == "" {
		//コマンド説明
		showHelp()
	} else {
		core.MsgErr("Invalid args!\n")
		showHelp()
		os.Exit(1)
	}

}
