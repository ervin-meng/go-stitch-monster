package config

type ConsulConfig struct {
	IP   string `mapstructure:"ip" json:"ip"`
	Port int    `mapstructure:"port" json:"port"`
}
