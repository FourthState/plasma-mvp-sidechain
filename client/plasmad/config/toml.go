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

# File containing unencrypted private key
# Used to sign eth transactions interacting with the rootchain  
ethereum_privkey_file = "{{ .EthPrivKeyFile }}"

# Gas limit for eth transactions
gas_limit = "{{.EthGasLimit }}"

# Boolean specifying if this node is a validator
is_validator = "{{ .IsValidator }}"`

var configTemplate *template.Template

func init() {
	var err error
	tmpl := template.New("plasmaConfigFileTemplate")
	if configTemplate, err = tmpl.Parse(defaultConfigTemplate); err != nil {
		panic(err)
	}
}

// parses the plasma.toml file and unmarshals it into a Config struct
func ParseConfig() (*Config, error) {
	config := DefaultConfig()
	err := viper.Unmarshal(config)
	return config, err
}

// WriteConfigFile renders config using the template and writes it to configFilePath.
func WriteConfigFile(configFilePath string, config *Config) {
	var buffer bytes.Buffer

	if err := configTemplate.Execute(&buffer, config); err != nil {
		panic(err)
	}

	// 0600 for owner only read+write permissions
	cmn.MustWriteFile(configFilePath, buffer.Bytes(), 0600)
}
