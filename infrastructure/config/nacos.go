package config

type NacosConfig struct {
	IP          string `mapstructure:"ip" json:"ip"`
	Port        int    `mapstructure:"port" json:"port"`
	NamespaceId string `mapstructure:"namespaceId" json:"namespaceId"`
	DataId      string `mapstructure:"dataId" json:"dataId"`
	Group       string `mapstructure:"group" json:"group"`
}
