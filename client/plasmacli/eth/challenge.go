package eth

/*
import (
	"fmt"
	"strconv"

	ks "github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	ethCmd.AddCommand(challengeCmd)
	challengeCmd.StringP(gasLimitF, "g", "21000", "gas limit for ethereum transaction")
	viper.BindPFlags(challengeCmd.Flags())
}

var challengeCmd = &cobra.Command{
	Use:   "challenge <exiting position> <challenging position> <account>",
	Short: "Challenge an existing exit",
	Long:  ``,
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		// parse positions
		exitingPos, err := plasma.FromPositionsString(args[0])
		if err != nil {
			return err
		}

		challengingPos, err := plasma.FromPositionsString(args[1])
		if err != nil {
			return err
		}

		gasLimit, err := strconv.ParseUint(viper.GetString(gasLimitF), 10, 64)
		if err != nil {
			return fmt.Error("failed to parse gas limit: { %s }", err)
		}

		key, err := ks.GetKey(args[2])
		if err != nil {
			return fmt.Errof("failed to retrieve account key: { %s }", err)
		}

		// bind key
		auth := bind.NewKeyedTransactor(key)
		transactOpts = bind.TransactOpts{
			From:     auth.From,
			Signer:   auth.Signer,
			GasLimit: gasLimit,
		}

		return nil
	},
}*/
