package cmd

import (
	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteKeyCmd)
}

var deleteKeyCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete the given key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		buf := client.BufferStdin()
		oldpass, err := client.GetPassword("DANGER - enter password to permanetly delete key:", buf)
		if err != nil {
			return err
		}

		kb, err := client.GetKeyBase()
		if err != nil {
			return err
		}

		err = kb.Delete(name, oldpass)
		if err != nil {
			return err
		}
		fmt.Println("Password deleted forever")
		return nil
	},
}
