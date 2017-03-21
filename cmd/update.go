package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/skuid/helm-value-store/dynamo"
	"github.com/skuid/spec"
	"github.com/spf13/cobra"
)

type updateCmdArgs struct {
	table string
	uuid  string

	file    string
	values  []string
	labels  spec.SelectorSet
	version string
}

var updateArgs = &updateCmdArgs{}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update a release in the release store",
	Long:  "Update a release. Any specified fields (other than UUID) will overwrite the existing fields.",
	Run:   update,
}

func init() {
	RootCmd.AddCommand(updateCmd)
	f := updateCmd.Flags()
	f.StringVar(&updateArgs.table, "table", "helm-charts", "Name of table")
	f.StringVar(&updateArgs.uuid, "uuid", "", "The UUID of the release")
	f.StringVarP(&updateArgs.file, "file", "f", "", "Name of values file. All values for this release will be overwritten with these values.")
	f.StringArrayVar(&updateArgs.values, "set", []string{}, `set values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2).
    	These values will be merged with the existing values in the store or with the values file provided by --file.`)
	f.VarP(&updateArgs.labels, "labels", "l", `The labels to apply. Each label should have the format "k=v".
    	Can be specified multiple times, or a comma-separated list.`)
	f.StringVar(&updateArgs.version, "version", "", "Version of the release")

	updateCmd.MarkFlagRequired("uuid")
	err := updateCmd.MarkFlagFilename("file", valueExtensions...)
	if err != nil {
		exitOnErr(err)
	}
}

func update(cmd *cobra.Command, args []string) {

	rs, err := dynamo.NewReleaseStore(updateArgs.table)
	exitOnErr(err)

	if len(updateArgs.uuid) == 0 {
		exitOnErr(errors.New("Must supply a UUID!"))
	}
	release, err := rs.Get(updateArgs.uuid)
	exitOnErr(err)

	if len(updateArgs.file) > 0 {
		values, err := ioutil.ReadFile(updateArgs.file)
		exitOnErr(err)
		release.Values = string(values)
	}

	if len(updateArgs.values) > 0 {
		err := release.MergeValues(updateArgs.values)
		exitOnErr(err)
	}

	if len(updateArgs.labels) > 0 {
		release.Labels = updateArgs.labels.ToMap()
	}
	if len(updateArgs.version) > 0 {
		release.Version = updateArgs.version
	}

	err = rs.Put(*release)
	exitOnErr(err)
	fmt.Printf("Updated release %s in release store!\n", release.Name)
}
