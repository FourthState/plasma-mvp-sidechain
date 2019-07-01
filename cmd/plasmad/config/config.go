package config

import (
	"bytes"
	"github.com/spf13/viper"
	cmn "github.com/tendermint/tendermint/libs/common"
	"text/template"
)

const defaultConfigTemplate = `# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

##### ethereum configuration #####

# Ethereum plasma contract address
ethereum_plasma_contract_address = "{{ .EthPlasmaContractAddr }}"

# Node URL for eth client
ethereum_nodeurl = "{{ .EthNodeURL }}"

# Number of Ethereum blocks until a submitted block header is considered final
ethereum_finality = "{{ .EthBlockFinality }}"

##### plasma configuration #####

# Plasma block commitment rate. i.e 1m30s, 1m, 1h, etc.
block_commitment_rate = "{{ .PlasmaCommitmentRate }}"

# Boolean specifying if this node is the operator of the plasma contract
is_operator = "{{ .IsOperator }}"

# Hex encoded private key
# Used to sign eth transactions interacting with the contract
operator_privatekey = "{{ .OperatorPrivateKey }}"`

// Must match the above defaultConfigTemplate
type PlasmaConfig struct {
	EthPlasmaContractAddr string `mapstructure:"ethereum_plasma_contract_address"`
	EthNodeURL            string `mapstructure:"ethereum_nodeurl"`
	EthBlockFinality      string `mapstructure:"ethereum_finality"`

	IsOperator           bool   `mapstructure:"is_operator"`
	OperatorPrivateKey   string `mapstructure:"operator_privatekey"`
	PlasmaCommitmentRate string `mapstructure:"block_commitment_rate"`
}

var configTemplate *template.Template

func init() {
	var err error
	tmpl := template.New("plasmaConfigFileTemplate")
	if configTemplate, err = tmpl.Parse(defaultConfigTemplate); err != nil {
		panic(err)
	}
}

func DefaultPlasmaConfig() PlasmaConfig {
	return PlasmaConfig{
		EthPlasmaContractAddr: "",
		EthNodeURL:            "http://localhost:8545",
		EthBlockFinality:      "16",

		IsOperator:           false,
		OperatorPrivateKey:   "",
		PlasmaCommitmentRate: "1m",
	}
}

// TestPlasmaConfig writes the plasma.toml file used for testing
// NodeURL powered by ganache locally
// Contract address and private key generated deterministically using the "plasma" moniker with ganache
func TestPlasmaConfig() PlasmaConfig {
	return PlasmaConfig{
		EthPlasmaContractAddr: "31E491FC70cDb231774c61B7F46d94699dacE664",
		EthNodeURL:            "http://localhost:8545",
		EthBlockFinality:      "0",

		IsOperator:           true,
		OperatorPrivateKey:   "9cd69f009ac86203e54ec50e3686de95ff6126d3b30a19f926a0fe9323c17181",
		PlasmaCommitmentRate: "1m",
	}
}

// parses the plasma.toml file and unmarshals it into a Config struct
func ParsePlasmaConfigFromViper() (PlasmaConfig, error) {
	config := DefaultPlasmaConfig()
	err := viper.Unmarshal(&config)
	return config, err
}

// WriteConfigFile renders config using the template and writes it to configFilePath.
func WritePlasmaConfigFile(configFilePath string, config PlasmaConfig) {
	var buffer bytes.Buffer

	if err := configTemplate.Execute(&buffer, &config); err != nil {
		panic(err)
	}

	// 0600 for owner only read+write permissions
	cmn.MustWriteFile(configFilePath, buffer.Bytes(), 0600)
}
