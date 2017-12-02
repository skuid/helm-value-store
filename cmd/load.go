package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/skuid/helm-value-store/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type loadCmdArgs struct {
	file  string
	setup bool
}

var loadArgs = &loadCmdArgs{}

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "load a json file of releases",
	Run:   load,
}

func init() {
	RootCmd.AddCommand(loadCmd)
	f := loadCmd.Flags()
	f.StringVar(&loadArgs.file, "file", "", "Name of file to ingest")
	f.BoolVar(&loadArgs.setup, "setup", false, "Setup the value store (may create resources).")

	loadCmd.MarkFlagRequired("file")
	err := loadCmd.MarkFlagFilename("file", valueExtensions...)
	if err != nil {
		exitOnErr(err)
	}
}

func load(cmd *cobra.Command, args []string) {
	fmt.Printf("Opening %s\n", loadArgs.file)

	f, err := os.Open(loadArgs.file)
	exitOnErr(err)
	defer f.Close()

	releases := []store.Release{}
	err = json.NewDecoder(f).Decode(&releases)
	exitOnErr(err)

	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("timeout"))
	defer cancel()
	if loadArgs.setup {
		err = releaseStore.Setup(ctx)
		exitOnErr(err)
	}

	err = releaseStore.Load(ctx, releases)
	exitOnErr(err)
	fmt.Printf("Loaded %d resources into %s\n", len(releases), viper.GetString("backend"))
}
