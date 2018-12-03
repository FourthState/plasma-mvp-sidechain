package config

// Defines basic configuration needed to interact with a rootchain contract
type Config struct {
	IsValidator    bool   `mapstructure:"is_validator"`
	EthPrivKeyFile string `mapstructure:"eth_privkey_file"`
	EthRootchain   string `mapstructure:"eth_rootchain"`
	EthNodeURL     string `mapstructure:"eth_nodeurl"`
	EthMinFees     string `mapstructure:"eth_min_fees"`
}

func DefaultConfig() *Config {
	return &Config{false, "", "", "", "0"}
}
