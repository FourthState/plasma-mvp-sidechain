package config

// Defines basic configuration needed to interact with a rootchain contract
type Config struct {
	IsValidator      bool
	EthPrivKeyFile   string
	EthRootchain     string
	EthNodeURL       string
	EthMinFees       string
	EthBlockFinality string
}

func DefaultConfig() *Config {
	return &Config{false, "", "", "", "0", "0"}
}
