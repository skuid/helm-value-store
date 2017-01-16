package cmd

import (
	"fmt"
	"os"
	"text/template"

	"github.com/skuid/helm-value-store/dynamo"
	"github.com/spf13/cobra"
)

type installCmdArgs struct {
	timeout  int64
	dryRun   bool
	table    string
	selector selectorSet
}

var installArgs = installCmdArgs{}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install a release",
	Run:   install,
}

var installTmpl *template.Template

func init() {
	RootCmd.AddCommand(installCmd)
	installCmd.Flags().StringVar(&installArgs.table, "table", "helm-charts", "Name of table")
	installCmd.Flags().Int64Var(&installArgs.timeout, "timeout", 300, "time in seconds to wait for any individual kubernetes operation (like Jobs for hooks)")
	installCmd.Flags().BoolVar(&installArgs.dryRun, "dry-run", false, "simulate an upgrade")
	installCmd.Flags().VarP(&installArgs.selector, "selector", "s", `The selectors to use. Each selector should have the format "k=v".
		Can be specified multiple times, or a comma-separated list.`)

	var err error
	installTmpl, err = template.New("InstallCmd").Parse(
		"REGION={{.Labels.region}} ENV={{.Labels.environment}} helm install{{ if .Name }} --name {{.Name}}{{ end }}{{if .Namespace}} --namespace {{.Namespace}}{{end}}{{if .Version}} --version {{.Version}}{{end}} {{.Chart}}  \n",
	)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func install(cmd *cobra.Command, args []string) {
	fmt.Println("Installing releases:")
	fmt.Println("")

	rs, err := dynamo.NewReleaseStore(installArgs.table)
	exitOnErr(err)

	releases, err := rs.List(installArgs.selector.ToMap())
	exitOnErr(err)

	for _, r := range releases {
		installTmpl.Execute(os.Stdout, &r)
	}
}
