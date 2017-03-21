package cmd

import (
	"encoding/json"
	"os"

	"github.com/skuid/helm-value-store/dynamo"
	"github.com/skuid/spec"
	"github.com/spf13/cobra"
)

type dumpCmdArgs struct {
	table   string
	label   spec.SelectorSet
	verbose bool
}

var dumpArgs = &dumpCmdArgs{}

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "dump the JSON representation of releases",
	Run:   dump,
}

func init() {
	RootCmd.AddCommand(dumpCmd)
	f := dumpCmd.Flags()
	f.StringVar(&dumpArgs.table, "table", "helm-charts", "Name of table")
	f.VarP(&dumpArgs.label, "label", "l", `The labels to filter by. Each label should have the format "k=v".
    	Can be specified multiple times, or a comma-separated list.`)

	f.BoolVar(&dumpArgs.verbose, "v", false, "Pretty-print the JSON")
}

func dump(cmd *cobra.Command, args []string) {
	rs, err := dynamo.NewReleaseStore(dumpArgs.table)
	exitOnErr(err)

	releases, err := rs.List(dumpArgs.label.ToMap())
	exitOnErr(err)

	encoder := json.NewEncoder(os.Stdout)
	if dumpArgs.verbose {
		encoder.SetIndent("", "    ")
	}
	err = encoder.Encode(releases)
	exitOnErr(err)
}
