package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/skuid/helm-value-store/store"
	"github.com/spf13/cobra"
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
	f.StringVar(&loadArgs.file, "file", "dynamoReleases.json", "Name of file to ingest")
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

	if loadArgs.setup {
		releaseStore.Setup()
	}

	err = releaseStore.Load(releases)
	exitOnErr(err)
	fmt.Println("Loaded resources into dynamo!")
}
