package cmd

import (
	"errors"
	"fmt"

	"github.com/skuid/helm-value-store/dynamo"
	"github.com/spf13/cobra"
)

type deleteCmdArgs struct {
	table     string
	file      string
	labels    selectorSet
	name      string
	chart     string
	namespace string
	version   string
	uuid      string
}

var deleteArgs = &deleteCmdArgs{}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a release in the relase store",
	Run:   delete,
}

func init() {
	RootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVar(&deleteArgs.table, "table", "helm-charts", "Name of table")
	deleteCmd.Flags().StringVar(&getArgs.uuid, "uuid", "", "The UUID to delete")
}

func delete(cmd *cobra.Command, args []string) {
	if len(deleteArgs.uuid) == 0 {
		exitOnErr(errors.New("Must supply a UUID!"))
	}

	rs, err := dynamo.NewReleaseStore(deleteArgs.table)
	exitOnErr(err)

	err = rs.Delete(deleteArgs.uuid)
	exitOnErr(err)

	fmt.Printf("Deleted release %s in release store.\n", deleteArgs.uuid)
}
