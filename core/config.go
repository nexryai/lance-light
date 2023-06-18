package core

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Default DefaultConfig `yaml:"default"`
	Ports   []PortConfig  `yaml:"ports"`
	Router  RouterConfig  `yaml:"router"`
}

type DefaultConfig struct {
	AllowAllIn  bool `yaml:"allowAllIn"`
	AllowAllOut bool `yaml:"allowAllOut"`
	AllowAllFwd bool `yaml:"allowAllFwd"`
}

type PortConfig struct {
	Port                  int    `yaml:"port"`
	Proto                 string `yaml:"proto"`
	AllowIPs              string `yaml:"allowIPs"`
	AllowCountry          string `yaml:"allowCountry"`
	DenyFromCloudProviders bool   `yaml:"denyFromCloudProviders"`
	DenyFromAbuseIPs      bool   `yaml:"denyFromAbuseIPs"`
	DenyFromTorIPs        bool   `yaml:"denyFromTorIPs"`
	AllowInterface        string `yaml:"allowInterface"`
}

type RouterConfig struct {
	ConfigAsRouter bool   `yaml:"configAsRouter"`
	WANInterface   string `yaml:"wanInterface"`
	LANInterface   string `yaml:"lanInterface"`
	ForceDNS       string `yaml:"forceDNS"`
}

func LoadConfig(configFilePath string) Config {
	// ファイルの読み込み
	data, err := ioutil.ReadFile(configFilePath)
	ExitOnError(err)

	// ファイルの内容を構造体にマッピング
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	// 読み込んだ構造体の内容を表示
	fmt.Printf("%+v\n", config)

	// portsの1番目の項目のallowIPsを取得
	if len(config.Ports) > 0 {
		allowIPs := config.Ports[0].AllowIPs
		fmt.Println("allowIPs:", allowIPs)
	}

	return config
}