package cmd

import (
	"fmt"

	"github.com/exoscale/egoscale"

	"github.com/spf13/cobra"
)

func init() {
	firewallDeleteCmd.Flags().BoolP("force", "f", false, "Attempt to remove security group without prompting for confirmation")
	firewallCmd.AddCommand(firewallDeleteCmd)
}

// deleteCmd represents the delete command
var firewallDeleteCmd = &cobra.Command{
	Use:     "delete <security group name | id>",
	Short:   "Delete security group",
	Aliases: gDeleteAlias,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Usage()
		}

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			return err
		}

		if !force {
			if !askQuestion(fmt.Sprintf("sure you want to delete %q security group", args[0])) {
				return nil
			}
		}

		return deleteFirewall(args[0])
	},
}

func deleteFirewall(name string) error {
	sg, err := getSecurityGroupByNameOrID(cs, name)
	if err != nil {
		return err
	}

	if err := cs.Delete(&egoscale.SecurityGroup{Name: sg.Name, ID: sg.ID}); err != nil {
		return err
	}

	println(sg.ID)
	return nil
}
