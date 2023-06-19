package main

import (
	"flag"
	"fmt"
	"lance-light/core"
	"lance-light/render"
)

// 成功したらTrue、そうでなければFalseを返す
func writeRulesFromConfig(configFilePath string) bool {
	rules := render.GenRulesFromConfig(configFilePath)
	for _, item := range rules {
		fmt.Println(item)
	}
	return true
}

func main() {
	core.MsgInfo("LanceLight ver0.01")

	// ファイルパスを格納するための変数を定義
	var configFilePath string

	flag.StringVar(&configFilePath, "f", "", "Path of config.yml")
	flag.StringVar(&configFilePath, "file", "", "Path of config.yml")

	// コマンドライン引数の解析
	flag.Parse()

	// filePath の値を使用して何かしらの処理を行う
	if configFilePath == "" {
		configFilePath = "/etc/lance.yml"
	}

	writeRulesFromConfig(configFilePath)
}
