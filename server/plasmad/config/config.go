package config

// Defines basic configuration needed to interact with a rootchain contract
type PlasmaConfig struct {
	IsOperator            bool
	EthPrivKeyFile        string
	EthPlasmaContractAddr string
	EthNodeURL            string
	EthBlockFinality      string
}

func DefaultPlasmaConfig() *PlasmaConfig {
	return &PlasmaConfig{false, "", "", "", "0"}
}
