package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/google/uuid"
	"github.com/skuid/helm-value-store/store"
	"github.com/skuid/spec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type createCmdArgs struct {
	file      string
	labels    spec.SelectorSet
	name      string
	chart     string
	namespace string
	version   string
}

var createArgs = &createCmdArgs{}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a release in the release store",
	Run:   create,
}

func init() {
	RootCmd.AddCommand(createCmd)
	f := createCmd.Flags()
	f.StringVarP(&createArgs.file, "file", "f", "", "Name of values file")
	f.VarP(&createArgs.labels, "labels", "l", `The labels to apply. Each label should have the format "k=v".
    	Can be specified multiple times, or a comma-separated list.`)
	f.StringVar(&createArgs.name, "name", "", "Name of the release")
	f.StringVar(&createArgs.chart, "chart", "", "Chart of the release")
	f.StringVar(&createArgs.namespace, "namespace", "default", "Namespace of the release")
	f.StringVar(&createArgs.version, "version", "", "Version of the release")

	err := createCmd.MarkFlagRequired("chart")
	if err != nil {
		exitOnErr(err)
	}
	err = createCmd.MarkFlagFilename("file", valueExtensions...)
	if err != nil {
		exitOnErr(err)
	}
}

func create(cmd *cobra.Command, args []string) {
	r := store.Release{
		UniqueID:  uuid.New().String(),
		Labels:    createArgs.labels.ToMap(),
		Name:      createArgs.name,
		Chart:     createArgs.chart,
		Namespace: createArgs.namespace,
		Version:   createArgs.version,
	}
	fmt.Printf("%#v\n", r)
	fmt.Println(r)

	if len(createArgs.file) > 0 {
		values, err := ioutil.ReadFile(createArgs.file)
		exitOnErr(err)
		r.Values = string(values)
	}
	if len(createArgs.chart) == 0 {
		exitOnErr(errors.New("No chart provided"))
	}
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("timeout"))
	defer cancel()
	err := releaseStore.Put(ctx, r)
	exitOnErr(err)
	fmt.Println("Created release in release store!")
}
