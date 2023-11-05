package core

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Nftables  NftablesConfig `yaml:"nftables"`
	Default   DefaultConfig  `yaml:"default"`
	IpSet     []IpSetConfig  `yaml:"ipset"`
	Security  SecurityConfig `yaml:"security"`
	Ports     []PortConfig   `yaml:"ports"`
	Router    RouterConfig   `yaml:"router"`
	Nat       []NatConfig    `yaml:"nat"`
	Report    ReportConfig   `yaml:"report"`
	DebugMode bool           `yaml:"debugMode"`
}

type NftablesConfig struct {
	NftablesFilePath string `yaml:"configFilePath"`
	IpDefineFilePath string `yaml:"ipDefineFilePath"`
}

type IpSetConfig struct {
	Name string   `yaml:"name"`
	Ip   []string `yaml:"ip"`
}

type DefaultConfig struct {
	AllowAllIn    bool `yaml:"allowAllIn"`
	AllowAllOut   bool `yaml:"allowAllOut"`
	AllowAllFwd   bool `yaml:"allowAllFwd"`
	AllowPing     bool `yaml:"allowPing"`
	EnableIPv6    bool `yaml:"enableIPv6"`
	EnableLogging bool `yaml:"enableLogging"`
}

type SecurityConfig struct {
	AlwaysDenyIP              []string              `yaml:"alwaysDenyIP"`
	AlwaysDenyASN             []string              `yaml:"alwaysDenyASN"`
	AlwaysDenyAbuseIP         bool                  `yaml:"alwaysDenyAbuseIP"`
	AlwaysDenyTor             bool                  `yaml:"alwaysDenyTor"`
	DisablePortScanProtection bool                  `yaml:"disablePortScanProtection"`
	DisableIpFragmentsBlock   bool                  `yaml:"disableIpFragmentsBlock"`
	CloudProtection           CloudProtectionConfig `yaml:"cloudProtection"`
}

type CloudProtectionConfig struct {
	EnableCloudProtection bool `yaml:"enableCloudProtection"`
	BlockPublicProxy      bool `yaml:"blockPublicProxy"`
	BlockAbuseIP          bool `yaml:"blockAbuseIP"`
	BlockBulletproofIP    bool `yaml:"blockBulletproofIP"`
}

type PortConfig struct {
	Port           int    `yaml:"port"`
	Proto          string `yaml:"proto"`
	AllowIP        string `yaml:"allowIP"`
	AllowCountry   string `yaml:"allowCountry"`
	AllowInterface string `yaml:"allowInterface"`
}

type RouterConfig struct {
	ConfigAsRouter          bool                 `yaml:"configAsRouter"`
	WANInterface            string               `yaml:"wanInterface"`
	PrivateNetworkAddresses []string             `yaml:"privateNetworks"`
	LANInterfaces           []string             `yaml:"lanInterfaces"`
	ForceDNS                string               `yaml:"forceDNS"`
	CustomRoutes            []CustomRoutesConfig `yaml:"customRoutes"`
}

type CustomRoutesConfig struct {
	AllowIP        string `yaml:"allowIP"`
	AllowInterface string `yaml:"allowInterface"`
	AllowDST       string `yaml:"allowDST"`
}

type NatConfig struct {
	Interface string `yaml:"interface"`
	AllowIP   string `yaml:"allowIP"`
	DstIP     string `yaml:"dstIP"`
	DstPort   string `yaml:"dstPort"`
	Proto     string `yaml:"proto"`
	NatTo     string `yaml:"natTo"`
}

type ReportConfig struct {
	AbuseIpDbAPIKey string   `yaml:"abuseIpDbAPIKey"`
	TrustedIPs      []string `yaml:"trustedIPs"`
	ReportInterval  int      `yaml:"reportIntervalMinutes"`
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
