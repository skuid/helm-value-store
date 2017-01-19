package cmd

import (
	"encoding/json"
	"os"

	"github.com/skuid/helm-value-store/dynamo"
	"github.com/spf13/cobra"
)

type dumpCmdArgs struct {
	table    string
	selector selectorSet
	verbose  bool
}

var dumpArgs = &dumpCmdArgs{}

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "dump the JSON representation of releases",
	Run:   dump,
}

func init() {
	RootCmd.AddCommand(dumpCmd)
	dumpCmd.Flags().StringVar(&dumpArgs.table, "table", "helm-charts", "Name of table")
	dumpCmd.Flags().VarP(&dumpArgs.selector, "selector", "s", `The selectors to use. Each selector should have the format "k=v".
    	Can be specified multiple times, or a comma-separated list.`)

	dumpCmd.Flags().BoolVar(&dumpArgs.verbose, "v", false, "Pretty-print the JSON")
}

func dump(cmd *cobra.Command, args []string) {
	rs, err := dynamo.NewReleaseStore(dumpArgs.table)
	exitOnErr(err)

	releases, err := rs.List(dumpArgs.selector.ToMap())
	exitOnErr(err)

	encoder := json.NewEncoder(os.Stdout)
	if dumpArgs.verbose {
		encoder.SetIndent("", "    ")
	}
	err = encoder.Encode(releases)
	exitOnErr(err)
}
