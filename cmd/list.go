package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/skuid/helm-value-store/dynamo"
	"github.com/spf13/cobra"
	"github.com/cloudfoundry/bytefmt"
)

type listCmdArgs struct {
	table    string
	selector selectorSet
}

var listArgs = &listCmdArgs{}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list the releases",
	Run:   list,
}

func init() {
	RootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVar(&listArgs.table, "table", "helm-charts", "Name of table")
	listCmd.Flags().VarP(&listArgs.selector, "selector", "s", `The selectors to use. Each selector should have the format "k=v".
    	Can be specified multiple times, or a comma-separated list.`)
}

func list(cmd *cobra.Command, args []string) {
	rs, err := dynamo.NewReleaseStore(listArgs.table)
	exitOnErr(err)

	releases, err := rs.List(listArgs.selector.ToMap())
	exitOnErr(err)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	columns := []string{
		"UniqueId", "Name", "Namespace", "Chart", "Version", "Labels", "Values",
	}

	fmt.Fprintln(w, strings.Join(columns, "\t"))

	for _, release := range releases {
		columns := []string{
			release.UniqueID,
			release.Name,
			release.Namespace,
			release.Chart,
			release.Version,
			fmt.Sprintf("%s", release.Labels),
			bytefmt.ByteSize(uint64(len(release.Values))),
		}
		fmt.Fprintln(w, strings.Join(columns, "\t"))
	}
	w.Flush()
}
