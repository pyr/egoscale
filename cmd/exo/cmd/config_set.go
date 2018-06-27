package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setCmd represents the set command
var configSetCmd = &cobra.Command{
	Use:   "set <account name>",
	Short: "Set an account as default",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		if allAccount == nil {
			log.Fatalf("No accounts defined")
		}

		if !isAccountExist(args[0]) {
			log.Fatalf("Account %q doesn't exist", args[0])
		}

		viper.Set("defaultAccount", args[0])

		if err := addAccount(viper.ConfigFileUsed(), nil, false); err != nil {
			log.Fatal(err)
		}

		println("Default profile set to", args[0])

	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
}
