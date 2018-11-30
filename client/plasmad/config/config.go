package config

// Defines basic configuration needed to interact with a rootchain contract
type Config struct {
	EthPrivKeyFile string `mapstructure:"priv_key_file"`
	EthGasLimit    string `mapstructure:"gas_limit"`
	IsValidator    bool   `mapstructure:"is_validator"`
}

func DefaultConfig() *Config {
	return &Config{"", "0", false}
}
