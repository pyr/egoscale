package cmd

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/exoscale/egoscale/cmd/exo/table"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var templateListCmd = &cobra.Command{
	Use:     "list [keyword]",
	Short:   "List all available templates",
	Aliases: gListAlias,
	Run: func(cmd *cobra.Command, args []string) {

		keyword := ""
		if len(args) >= 1 {
			keyword = strings.Join(args, " ")
		}

		infos, err := listTemplates(keyword)
		if err != nil {
			log.Fatal(err)
		}

		table := table.NewTable(os.Stdout)
		table.SetHeader([]string{"Operating System", "Disk", "Release Date", "ID"})

		for _, v := range infos {
			sz := strconv.FormatInt(v.Size, 10)
			if sz == "0" {
				sz = ""
			}
			table.Append([]string{v.Name, sz, v.Created, v.ID})
		}
		table.Render()
	},
}

func init() {
	templateCmd.AddCommand(templateListCmd)
}
