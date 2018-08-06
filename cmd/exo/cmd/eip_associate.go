package cmd

import (
	"fmt"
	"net"

	"github.com/exoscale/egoscale"

	"github.com/spf13/cobra"
)

// associateCmd represents the associate command
var eipAssociateCmd = &cobra.Command{
	Use:     "associate <IP address> <instance name | instance id> [instance name | instance id] [...]",
	Short:   "Associate an IP to instance(s)",
	Aliases: gAssociateAlias,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return cmd.Usage()
		}

		for _, arg := range args[1:] {
			res, err := associateIP(args[0], arg)
			if err != nil {
				return err
			}
			println(res)
		}
		return nil
	},
}

func associateIP(ipAddr, instance string) (string, error) {
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return "", fmt.Errorf("invalid IP address %q", ipAddr)
	}

	vm, err := getVMWithNameOrID(instance)
	if err != nil {
		return "", err
	}

	defaultNic := vm.DefaultNic()

	if defaultNic == nil {
		return "", fmt.Errorf("the instance %q has not default NIC", vm.ID)
	}

	resp, err := cs.Request(&egoscale.AddIPToNic{NicID: defaultNic.ID, IPAddress: ip})
	if err != nil {
		return "", err
	}

	result := resp.(*egoscale.NicSecondaryIP)

	return result.NicID, nil
}

func init() {
	eipCmd.AddCommand(eipAssociateCmd)
}
