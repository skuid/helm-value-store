package cmd

import (
	"fmt"
	"strings"

	"github.com/skuid/helm-value-store/dynamo"
	"github.com/skuid/helm-value-store/store"
	"github.com/skuid/spec"
	"github.com/spf13/cobra"
)

type installCmdArgs struct {
	timeout int64
	dryRun  bool
	table   string
	labels  spec.SelectorSet
	values  []string

	uuid string
	name string
}

var installArgs = installCmdArgs{}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install or upgrade a release",
	Long:  "Install a release in the cluster and use the values from the value store",
	Run:   install,
}

func init() {
	RootCmd.AddCommand(installCmd)
	f := installCmd.Flags()
	f.StringVar(&installArgs.table, "table", "helm-charts", "Name of table")
	f.Int64Var(&installArgs.timeout, "timeout", 300, "time in seconds to wait for any individual kubernetes operation (like Jobs for hooks)")
	f.BoolVar(&installArgs.dryRun, "dry-run", false, "simulate an install/upgrade")
	f.VarP(&installArgs.labels, "label", "l", `The labels to filter by. Each label should have the format "k=v".
		Can be specified multiple times, or a comma-separated list.`)
	f.StringVar(&installArgs.uuid, "uuid", "", "The UUID to install. Takes precedence over --name")
	f.StringVar(&installArgs.name, "name", "", `The name of the release to install. If multiple releases of the same name are found,
		the install will fail. Use selectors to pair down releases`)
	f.StringArrayVar(&installArgs.values, "set", []string{}, "set values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
}

func releasesByName(name string, releases store.Releases) (release store.Releases) {
	response := store.Releases{}
	for _, r := range releases {
		if r.Name == name {
			response = append(response, r)
		}
	}
	return response
}

func install(cmd *cobra.Command, args []string) {
	rs, err := dynamo.NewReleaseStore(installArgs.table)
	exitOnErr(err)

	release := &store.Release{}

	if len(installArgs.uuid) > 0 {
		release, err = rs.Get(installArgs.uuid)
		exitOnErr(err)
	} else if len(installArgs.name) > 0 {
		releases, err := rs.List(installArgs.labels.ToMap())
		exitOnErr(err)

		matches := releasesByName(installArgs.name, releases)
		if len(matches) > 1 {
			exitOnErr(fmt.Errorf("Too many releases by the name: %s", installArgs.name))
		} else if len(matches) < 1 {
			exitOnErr(fmt.Errorf("No releases by the name: %s", installArgs.name))
		}
		release = &matches[0]
	} else {
		exitOnErr(fmt.Errorf("No release specified! Use %s or %s", "--name", "--uuid"))
	}
	_, getErr := release.Get()

	if getErr != nil && !strings.Contains(getErr.Error(), "not found") {
		exitOnErr(err)
	}

	if len(installArgs.values) > 0 {
		err := release.MergeValues(installArgs.values)
		exitOnErr(err)
	}

	dlLocation, err := release.Download()
	exitOnErr(err)
	fmt.Printf("Fetched chart %s to %s\n", release.Chart, dlLocation)

	if getErr != nil && strings.Contains(getErr.Error(), "not found") {
		// Install
		fmt.Printf("Installing Release %s\n", release)

		response, err := release.Install(dlLocation, installArgs.dryRun, installArgs.timeout)
		exitOnErr(err)
		fmt.Printf("Successfully installed release %s!\n", response.Release.Name)
	} else if getErr == nil {
		// Update
		fmt.Printf("Updating Release %s\n", release)
		response, err := release.Upgrade(dlLocation, installArgs.dryRun, installArgs.timeout)
		exitOnErr(err)
		fmt.Printf("Successfully upgraded release %s!\n", response.Release.Name)
	}

}
