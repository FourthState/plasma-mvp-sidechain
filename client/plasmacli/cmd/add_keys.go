package cmd

import (
	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/spf13/cobra"
)

const (
	flagRecover  = "recover"
	flagNoBackup = "no-backup"
)

func init() {
	addKeyCmd.Flags().Bool(flagRecover, false, "Provide seed phrase to recover existing key instead of creating a new key")
	addKeyCmd.Flags().Bool(flagNoBackup, false, "Do not print out seed phrase (others may be watching the terminal)")
	rootCmd.AddCommand(addKeyCmd)
}

var addKeyCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Create a new key, or import from seed",
	Long:  `Add a public/private key pair to the key manager.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		buf := client.BufferStdin()
		if len(args) != 1 || len(args[0]) == 0 {
			return errors.New("You must provide a name for the key")
		}
		name := args[0]
		kb, err := client.GetKeyBase()
		if err != nil {
			return err
		}

		_, err := kb.Get(name)
		if err == nil {
			return errors.New("A key already exists with the provided name")
		}

		pass, err := client.GetCheckPassword("Enter a passphrase for your key:", "Repeat the passphrase:", buf)
		if err != nil {
			return err
		}

		if viper.GetBool(flagRecover) {
			seed, err := client.GetSeed("Enter your recovery seed phrase:", buf)
			if err != nil {
				return err
			}
			info, err := kb.Recover(name, pass, seed)
			if err != nil {
				return err
			}
			viper.Set(flagNoBackup, true)
			client.PrintInfo(info)
		} else {
			info, seed, err := kb.Create(name, pass, nil)
			if err != nil {
				return err
			}
			client.PrintInfo(info)
			fmt.Println("**Important** write this seed phrase in a safe place.")
			fmt.Println("It is the only way to recover your account if you ever forget your password")
			fmt.Println()
			fmt.Println(seed)
		}
	},
}
