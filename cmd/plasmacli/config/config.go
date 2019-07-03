package config

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	cmn "github.com/tendermint/tendermint/libs/common"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const defaultConfigTemplate = `# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

##### ethereum configuration #####

# Ethereum plasma contract address
ethereum_contract_address = "{{ .EthPlasmaContractAddr }}"

# Node URL for eth client
ethereum_nodeurl = "{{ .EthNodeURL }}"


##### plasamd configuration #####

# Node URL for plasmad
node = "{{ .PlasmadNodeURL }}"

# Trust the connected plasmad node (don't verify proofs for responses)
trust_node = {{ .PlasmadTrustNode }}

# Chain identifier. Must be set if trust-node == false
chain_id = "{{ .PlasmadChainID }}"`

// Must match the above defaultConfigTemplate
type Config struct {
	// Ethereum config
	EthPlasmaContractAddr string `mapstructure:"ethereum_contract_address"`
	EthNodeURL            string `mapstructure:"ethereum_nodeurl"`

	// Plasmad config
	PlasmadNodeURL   string `mapstructure:"node"`
	PlasmadTrustNode bool   `mapstructure:"trust_node"`
	PlasmadChainID   string `mapstructure:"chain_id"`
}

var configTemplate *template.Template

func init() {
	var err error
	tmpl := template.New("configFileTemplate")
	if configTemplate, err = tmpl.Parse(defaultConfigTemplate); err != nil {
		panic(err)
	}
}

func DefaultConfig() Config {
	return Config{
		EthPlasmaContractAddr: "",
		EthNodeURL:            "http://localhost:8545",
		PlasmadNodeURL:        "tcp://localhost:26657",
		PlasmadTrustNode:      false,
		PlasmadChainID:        "",
	}
}

// RegisterViper will match client flags with config and register env
func RegisterViperAndEnv() {
	viper.SetEnvPrefix("PCLI")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// The configuration files use underscores while the Cosmos SDK uses
	// hypens. These aliases align `Viper.Get(..)` for both
	// the SDK and the configuration file
	viper.RegisterAlias("trust_node", "trust-node")
	viper.RegisterAlias("chain_id", "chain-id")
}

// parses the plasma.toml file and unmarshals it into a Config struct
func ParseConfigFromViper() (Config, error) {
	config := Config{}
	err := viper.Unmarshal(&config)
	return config, err
}

// WriteConfigFile renders config using the template and writes it to configFilePath.
func WriteConfigFile(configFilePath string, config Config) error {
	var buffer bytes.Buffer

	if err := configTemplate.Execute(&buffer, &config); err != nil {
		return fmt.Errorf("template: %s", err)
	}

	if err := cmn.EnsureDir(filepath.Dir(configFilePath), os.ModePerm); err != nil {
		return fmt.Errorf("ensuredir: %s", err)
	}

	// 0666 allows for read and write for any user
	if err := cmn.WriteFile(configFilePath, buffer.Bytes(), 0666); err != nil {
		return fmt.Errorf("writefile: %s", err)
	}

	return nil
}
