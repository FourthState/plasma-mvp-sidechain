package config

import (
	"bytes"
	"text/template"

	"github.com/spf13/viper"
	cmn "github.com/tendermint/tendermint/libs/common"
)

const defaultConfigTemplate = `# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

##### ethereum config options #####
# Boolean specifying if this node is the operator of the plasma contract
is_operator = "{{ .IsOperator }}"

# File containing unencrypted private key
# Used to sign eth transactions interacting with the rootchain 
# Default directory is $HOME/.plasmad/config/
ethereum_operator_privatekey = "{{ .EthOperatorPrivateKey }}"

# Ethereum rootchain contract address
ethereum_plasma_contract_address = "{{ .EthPlasmaContractAddr }}"

# Node URL for eth client
ethereum_nodeurl = "{{ .EthNodeURL }}"

# Number of Ethereum blocks until a submitted block header is considered final
ethereum_finality = "{{ .EthBlockFinality }}"`

var configTemplate *template.Template

func init() {
	var err error
	tmpl := template.New("plasmaConfigFileTemplate")
	if configTemplate, err = tmpl.Parse(defaultConfigTemplate); err != nil {
		panic(err)
	}
}

// parses the plasma.toml file and unmarshals it into a Config struct
func ParsePlasmaConfig() (*PlasmaConfig, error) {
	config := DefaultPlasmaConfig()
	err := viper.Unmarshal(config)
	return config, err
}

// WriteConfigFile renders config using the template and writes it to configFilePath.
func WritePlasmaConfigFile(configFilePath string, config *PlasmaConfig) {
	var buffer bytes.Buffer

	if err := configTemplate.Execute(&buffer, config); err != nil {
		panic(err)
	}

	// 0600 for owner only read+write permissions
	cmn.MustWriteFile(configFilePath, buffer.Bytes(), 0600)
}
