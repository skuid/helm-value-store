package cmd

import (
	"errors"
	"fmt"

	"github.com/skuid/helm-value-store/dynamo"
	"github.com/spf13/cobra"
)

type getCmdArgs struct {
	timeout  int64
	dryRun   bool
	table    string
	selector selectorSet

	uuid string
	name string
}

var getArgs = getCmdArgs{}

var getCmd = &cobra.Command{
	Use:   "get-values",
	Short: "get the values of a release",
	Run:   get,
}

func init() {
	RootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVar(&getArgs.table, "table", "helm-charts", "Name of table")
	getCmd.Flags().StringVar(&getArgs.uuid, "uuid", "", "The UUID to get.")

	getCmd.MarkFlagRequired("uuid")
}

func get(cmd *cobra.Command, args []string) {
	rs, err := dynamo.NewReleaseStore(getArgs.table)
	exitOnErr(err)

	if len(getArgs.uuid) == 0 {
		exitOnErr(errors.New("Must supply a UUID!"))
	}
	release, err := rs.Get(getArgs.uuid)
	exitOnErr(err)

	fmt.Print(release.Values)
}
