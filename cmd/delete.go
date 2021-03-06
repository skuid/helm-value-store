package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/skuid/spec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type deleteCmdArgs struct {
	file      string
	labels    spec.SelectorSet
	name      string
	chart     string
	namespace string
	version   string
	uuid      string
}

var deleteArgs = &deleteCmdArgs{}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a release in the release store",
	Run:   delete,
}

func init() {
	RootCmd.AddCommand(deleteCmd)
	f := deleteCmd.Flags()
	f.StringVar(&deleteArgs.uuid, "uuid", "", "The UUID to delete")

	deleteCmd.MarkFlagRequired("uuid")
}

func delete(cmd *cobra.Command, args []string) {
	if len(deleteArgs.uuid) == 0 {
		exitOnErr(errors.New("Must supply a UUID"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("timeout"))
	defer cancel()
	err := releaseStore.Delete(ctx, deleteArgs.uuid)
	exitOnErr(err)

	fmt.Printf("Deleted release %s in release store.\n", deleteArgs.uuid)
}
