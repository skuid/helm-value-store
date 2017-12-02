package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"code.cloudfoundry.org/bytefmt"
	"github.com/skuid/helm-value-store/store"
	"github.com/skuid/spec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type listCmdArgs struct {
	labels spec.SelectorSet
	name   string
}

var listArgs = &listCmdArgs{}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list the releases",
	Run:   list,
}

func init() {
	RootCmd.AddCommand(listCmd)
	f := listCmd.Flags()
	f.VarP(&listArgs.labels, "labels", "l", `The labels to filter by. Each label should have the format "k=v".
    	Can be specified multiple times, or a comma-separated list.`)
	f.StringVar(&listArgs.name, "name", "", "Filter by release name")
}

func filterByName(releases store.Releases, name string) store.Releases {
	response := store.Releases{}
	for _, r := range releases {
		if r.Name == name {
			response = append(response, r)
		}
	}
	return response
}

func list(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("timeout"))
	defer cancel()
	releases, err := releaseStore.List(ctx, listArgs.labels.ToMap())
	exitOnErr(err)

	if len(listArgs.name) > 0 {
		releases = filterByName(releases, listArgs.name)
	}

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
