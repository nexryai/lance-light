package main

import (
	"flag"
	"fmt"
	"lance-light/internal/config"
	"lance-light/internal/log"
	"lance-light/internal/render"
	"lance-light/internal/system"
	"os"
)

func applyNftablesRules(configFilePath string) {
	system.ExecCommand("nft", []string{"-f", configFilePath})
}

func flushNftablesRules() {
	system.ExecCommand("nft", []string{"flush", "table", "inet", "lance"})
}

func checkConfigFile(path string) {
	system.ExecCommand("nft", []string{"--check", "-f", path})
}

func writeRulesFromConfig(config *config.Config) bool {
	// ipdefine.conf (IPのリストを定義するやつ)を生成
	ipDefineRules, err := render.GenIpDefineRules(config)
	if err != nil {
		log.ExitOnError(err, "Network Error. Please use offline mode!")
	} else {
		system.WriteToFile(ipDefineRules, config.Nftables.IpDefineFilePath)
	}

	// nftablesルールを生成
	rules := render.GenRulesFromConfig(config)
	system.WriteToFile(rules, config.Nftables.NftablesFilePath)
	return true
}

func exportRulesFromConfig(config *config.Config) bool {
	ipDefineRules, _ := render.GenIpDefineRules(config)
	for _, item := range ipDefineRules {
		fmt.Println(item)
	}

	rules := render.GenRulesFromConfig(config)
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
		"-f [PATH]: Specify the path to the configuration file (Default: /etc/lance.yml)")
}

func main() {
	configFilePath := flag.String("f", "/etc/lance.yml", "Path of config.yml")
	debugMode := flag.Bool("debug", false, "Enable debug mode")

	flag.Parse()

	// 現状デバッグモードはログ以外に影響を与えないが将来的に変わる可能性もある。
	// 環境変数以外にもいいやり方あるかもしれない。毎回configをcore.MsgDbg呼ぶ時に引数渡すのは面倒だから避けたい。
	if *debugMode {
		log.MsgInfo("debug mode!")
		err := os.Setenv("LANCE_DEBUG_MODE", "true")
		log.ExitOnError(err, "Failed to set environment variable.")
	} else {
		err := os.Setenv("LANCE_DEBUG_MODE", "false")
		log.ExitOnError(err, "Failed to set environment variable.")
	}

	log.MsgDebug("configFilePath: " + *configFilePath)

	config := config.LoadConfig(*configFilePath)

	if config.Nftables.NftablesFilePath == "" {
		config.Nftables.NftablesFilePath = "/etc/nftables.lance.conf"
	}

	if config.Nftables.IpDefineFilePath == "" {
		config.Nftables.IpDefineFilePath = "/etc/nftables.ipdefine.conf"
	}

	operation := flag.Arg(0)

	if operation == "apply" {
		writeRulesFromConfig(&config)
		checkConfigFile(config.Nftables.NftablesFilePath)

		// nftコマンドを実行して適用
		flushNftablesRules()
		applyNftablesRules(config.Nftables.NftablesFilePath)

		log.MsgInfo("Firewall settings have been applied successfully.")
	} else if operation == "enable" {
		writeRulesFromConfig(&config)
		checkConfigFile(config.Nftables.NftablesFilePath)

		// nftコマンドを実行して適用
		applyNftablesRules(config.Nftables.NftablesFilePath)
		log.MsgInfo("LanceLight firewall is enabled.")
	} else if operation == "offline" {
		// Q.これは何
		// A.オフライン環境だとレンダリングできない（CloudflareのIPなどが取得できない）。起動直後などのオフラインな環境でも最低限の保護を有効にするため、一旦lance.ymlの変更を反映せずとりあえず古いルールをロードだけする。

		// nftコマンドを実行して適用
		applyNftablesRules(config.Nftables.NftablesFilePath)
		log.MsgInfo("LanceLight firewall is enabled. (Offline mode!)")
	} else if operation == "export" {
		// エクスポート
		log.MsgDebug(fmt.Sprintf("configFilePath: %s", *configFilePath))
		exportRulesFromConfig(&config)
	} else if operation == "disable" {
		// 設定をアンロードする
		flushNftablesRules()
		log.MsgInfo("LanceLight firewall is disabled.")
	} else if operation == "" {
		//コマンド説明
		showHelp()
	} else {
		log.MsgErr("Invalid args!\n")
		showHelp()
		os.Exit(1)
	}

}
