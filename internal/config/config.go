package config

import (
	"gopkg.in/yaml.v2"
	"lance-light/internal/log"
	"os"
)

type Config struct {
	Nftables  NftablesConfig `yaml:"nftables"`
	Default   DefaultConfig  `yaml:"default"`
	IpSet     []IpSetConfig  `yaml:"ipset"`
	Security  SecurityConfig `yaml:"security"`
	Ports     []PortConfig   `yaml:"ports"`
	Outgoing  OutgoingConfig `yaml:"outbound"`
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
	Url  string   `yaml:"url"`
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
	AlwaysDenyIP              []string `yaml:"alwaysDenyIP"`
	AlwaysDenyASN             []string `yaml:"alwaysDenyASN"`
	AlwaysDenyAbuseIP         bool     `yaml:"alwaysDenyAbuseIP"`
	AlwaysDenyTor             bool     `yaml:"alwaysDenyTor"`
	DisablePortScanProtection bool     `yaml:"disablePortScanProtection"`
	DisableIpFragmentsBlock   bool     `yaml:"disableIpFragmentsBlock"`
}

type PortConfig struct {
	Port           int    `yaml:"port"`
	Proto          string `yaml:"proto"`
	AllowIP        string `yaml:"allowIP"`
	AllowInterface string `yaml:"allowInterface"`
}

type OutgoingConfig struct {
	Compatibility []string              `yaml:"compatibility"`
	Allowed       []OutgoingAllowConfig `yaml:"allowed"`
}

type OutgoingAllowConfig struct {
	Dport string `yaml:"dport"`
	Proto string `yaml:"proto"`
	DstIP string `yaml:"dstIP"`
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
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		log.MsgFatalAndExit(err, "An error occurred while loading the configuration file. Are the configuration file paths and permissions correct?")
	}

	// ファイルの内容を構造体にマッピング
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.MsgFatalAndExit(err, "The configuration file was loaded successfully, but the mapping failed.")
	}

	return config
}
