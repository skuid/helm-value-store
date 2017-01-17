package cmd

import (
	"fmt"

	"github.com/skuid/helm-value-store/dynamo"
	"github.com/skuid/helm-value-store/store"
	"github.com/spf13/cobra"
)

type installCmdArgs struct {
	timeout  int64
	dryRun   bool
	table    string
	selector selectorSet

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
	installCmd.Flags().StringVar(&installArgs.table, "table", "helm-charts", "Name of table")
	installCmd.Flags().Int64Var(&installArgs.timeout, "timeout", 300, "time in seconds to wait for any individual kubernetes operation (like Jobs for hooks)")
	installCmd.Flags().BoolVar(&installArgs.dryRun, "dry-run", false, "simulate an install/upgrade")
	installCmd.Flags().VarP(&installArgs.selector, "selector", "s", `The selectors to use. Each selector should have the format "k=v".
		Can be specified multiple times, or a comma-separated list.`)
	installCmd.Flags().StringVar(&installArgs.uuid, "uuid", "", "The UUID to install. Takes precedence over --name")
	installCmd.Flags().StringVar(&installArgs.name, "name", "", `The name of the release to install. If multiple releases of the same name are found,
		the install will fail. Use selectors to pair down releases`)
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
		releases, err := rs.List(installArgs.selector.ToMap())
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

	if getErr != nil && getErr.Error() != "rpc error: code = 2 desc = release: not found" {
		exitOnErr(err)
	}

	dlLocation, err := release.Download()
	exitOnErr(err)
	fmt.Printf("Fetched chart %s to %s\n", release.Chart, dlLocation)

	if getErr != nil && getErr.Error() == "rpc error: code = 2 desc = release: not found" {
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
