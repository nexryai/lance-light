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

func writeRulesFromConfig(cfg *config.Config) bool {
	// ipdefine.conf (IPのリストを定義するやつ)を生成
	ipDefineRules, err := render.GenIpDefineRules(cfg)
	if err != nil {
		log.MsgFatalAndExit(err, "Network Error. Please use offline mode!")
	} else {
		system.WriteToFile(ipDefineRules, cfg.Nftables.IpDefineFilePath)
	}

	// nftablesルールを生成
	rules := render.GenRulesFromConfig(cfg)
	system.WriteToFile(rules, cfg.Nftables.NftablesFilePath)
	return true
}

func exportRulesFromConfig(cfg *config.Config) bool {
	ipDefineRules, _ := render.GenIpDefineRules(cfg)
	for _, item := range ipDefineRules {
		fmt.Println(item)
	}

	rules := render.GenRulesFromConfig(cfg)
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
	configFilePath := flag.String("f", "/etc/lance.yml", "Path of cfg.yml")
	debugMode := flag.Bool("debug", false, "Enable debug mode")

	flag.Parse()

	// 現状デバッグモードはログ以外に影響を与えないが将来的に変わる可能性もある。
	// 環境変数以外にもいいやり方あるかもしれない。毎回configをcore.MsgDbg呼ぶ時に引数渡すのは面倒だから避けたい。
	if *debugMode {
		log.MsgInfo("debug mode!")
		err := os.Setenv("LANCE_DEBUG_MODE", "true")
		if err != nil {
			log.MsgFatalAndExit(err, "Failed to set environment variable.")
		}
	} else {
		err := os.Setenv("LANCE_DEBUG_MODE", "false")
		if err != nil {
			log.MsgFatalAndExit(err, "Failed to set environment variable.")
		}
	}

	log.MsgDebug("configFilePath: " + *configFilePath)

	cfg := config.LoadConfig(*configFilePath)

	if cfg.Nftables.NftablesFilePath == "" {
		cfg.Nftables.NftablesFilePath = "/etc/nftables.lance.conf"
	}

	if cfg.Nftables.IpDefineFilePath == "" {
		cfg.Nftables.IpDefineFilePath = "/etc/nftables.ipdefine.conf"
	}

	operation := flag.Arg(0)

	if operation == "apply" {
		writeRulesFromConfig(&cfg)
		checkConfigFile(cfg.Nftables.NftablesFilePath)

		// nftコマンドを実行して適用
		flushNftablesRules()
		applyNftablesRules(cfg.Nftables.NftablesFilePath)

		log.MsgInfo("Firewall settings have been applied successfully.")
	} else if operation == "enable" {
		writeRulesFromConfig(&cfg)
		checkConfigFile(cfg.Nftables.NftablesFilePath)

		// nftコマンドを実行して適用
		applyNftablesRules(cfg.Nftables.NftablesFilePath)
		log.MsgInfo("LanceLight firewall is enabled.")
	} else if operation == "offline" {
		// Q.これは何
		// A.オフライン環境だとレンダリングできない（CloudflareのIPなどが取得できない）。起動直後などのオフラインな環境でも最低限の保護を有効にするため、一旦lance.ymlの変更を反映せずとりあえず古いルールをロードだけする。

		// nftコマンドを実行して適用
		applyNftablesRules(cfg.Nftables.NftablesFilePath)
		log.MsgInfo("LanceLight firewall is enabled. (Offline mode!)")
	} else if operation == "export" {
		// エクスポート
		log.MsgDebug(fmt.Sprintf("configFilePath: %s", *configFilePath))
		exportRulesFromConfig(&cfg)
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
