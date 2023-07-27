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

func writeRulesFromConfig(configFilePath string, nftablesFilePath string, ipDefineFilePath string, addFlushRule bool) bool {
	config := core.LoadConfig(configFilePath)

	ipDefineRules, err := render.GenIpDefineRules("cloudflare", &config)
	if err != nil {
		//ToDo: ipDefineFilePathが存在しなければ失敗扱い
	} else {
		core.WriteToFile(ipDefineRules, ipDefineFilePath)
	}

	rules := render.GenRulesFromConfig(&config, addFlushRule)
	core.WriteToFile(rules, nftablesFilePath)
	return true
}

func exportRulesFromConfig(configFilePath string) bool {
	config := core.LoadConfig(configFilePath)

	rules := render.GenRulesFromConfig(&config, false)
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
		"-f [PATH]: Specify the path to the configuration file (Default: /etc/lance.yml)\n",
		"-o [PATH]: Where to write nftables rules. Need not to use except for debugging. (Default: /etc/nftables.lance.conf)\n\n")
}

func main() {
	configFilePath := flag.String("f", "/etc/lance.yml", "Path of config.yml")
	nftablesFilePath := flag.String("o", "/etc/nftables.lance.conf", "Path of nftables.conf")
	ipDefineFilePath := flag.String("d", "/etc/nftables.ipdefine.conf", "Path of ipdefine.conf")
	debugMode := flag.Bool("debug", false, "Enable debug mode")

	flag.Parse()

	// 現状デバッグモードはログ以外に影響を与えないが将来的に変わる可能性もある。
	// 環境変数以外にもいいやり方あるかもしれない。毎回core.MsgDbg呼ぶ時に引数渡すのは面倒だから避けたい。
	if *debugMode {
		core.MsgInfo("debug mode!")
		os.Setenv("LANCE_DEBUG_MODE", "true")
	} else {
		os.Setenv("LANCE_DEBUG_MODE", "false")
	}

	operation := flag.Arg(0)

	if operation == "apply" {

		// 設定をリセットして再設定する
		writeRulesFromConfig(*configFilePath, *nftablesFilePath, *ipDefineFilePath, true)

		// nftコマンドを実行して適用
		applyNftablesRules(*nftablesFilePath)

		core.MsgInfo("Firewall settings have been applied successfully.")

	} else if operation == "enable" {

		// 設定を適用する
		writeRulesFromConfig(*configFilePath, *nftablesFilePath, *ipDefineFilePath, false)

		// nftコマンドを実行して適用
		applyNftablesRules(*nftablesFilePath)

		core.MsgInfo("LanceLight firewall is enabled.")

	} else if operation == "export" {

		// エクスポート
		core.MsgDebug(fmt.Sprintf("configFilePath: %s", *configFilePath))
		exportRulesFromConfig(*configFilePath)

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
