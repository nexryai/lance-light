package core

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Default  DefaultConfig  `yaml:"default"`
	Security SecurityConfig `yaml:"security"`
	Ports    []PortConfig   `yaml:"ports"`
	Router   RouterConfig   `yaml:"router"`
}

type DefaultConfig struct {
	AllowAllIn  bool `yaml:"allowAllIn"`
	AllowAllOut bool `yaml:"allowAllOut"`
	AllowAllFwd bool `yaml:"allowAllFwd"`
	AllowPing   bool `yaml:"allowPing"`
	EnableIPv6  bool `yaml:"enableIPv6"`
}

type SecurityConfig struct {
	AlwaysDenyIP      []string `yaml:"alwaysDenyIP"`
	AlwaysDenyASN     []string `yaml:"alwaysDenyASN"`
	AlwaysDenyAbuseIP bool     `yaml:"alwaysDenyAbuseIP"`
	AlwaysDenyTor     bool     `yaml:"alwaysDenyTor"`
}

type PortConfig struct {
	Port             int    `yaml:"port"`
	Proto            string `yaml:"proto"`
	AllowIP          string `yaml:"allowIP"`
	AllowCountry     string `yaml:"allowCountry"`
	DenyFromAbuseIPs bool   `yaml:"denyFromAbuseIPs"`
	DenyFromTorIPs   bool   `yaml:"denyFromTorIPs"`
	AllowInterface   string `yaml:"allowInterface"`
}

type RouterConfig struct {
	ConfigAsRouter          bool     `yaml:"configAsRouter"`
	WANInterface            string   `yaml:"wanInterface"`
	PrivateNetworkAddresses []string `yaml:"privateNetworks"`
	LANInterfaces           []string `yaml:"lanInterfaces"`
	ForceDNS                string   `yaml:"forceDNS"`
}

func LoadConfig(configFilePath string) Config {
	// ファイルの読み込み
	data, err := ioutil.ReadFile(configFilePath)
	ExitOnError(err, "An error occurred while loading the configuration file. Are the configuration file paths and permissions correct?")

	// ファイルの内容を構造体にマッピング
	var config Config
	err = yaml.Unmarshal(data, &config)
	ExitOnError(err, "The configuration file was loaded successfully, but the mapping failed.")

	/*
		if len(config.Ports) > 0 {
			allowIPs := config.Ports[0].AllowIPs
			fmt.Println("allowIPs:", allowIPs)
		}
	*/

	return config
}
