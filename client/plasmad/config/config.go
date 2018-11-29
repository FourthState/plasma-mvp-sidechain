package config

// Defines basic configuration needed to interact with a rootchain contract
type Config struct {
	EthPrivKeyFile string `mapstructure:"priv_key_file"`
	EthGasLimit    string `mapstructure:"gas_limit"`
}

func DefaultConfig() *Config {
	return &Config{"", "0"}
}
