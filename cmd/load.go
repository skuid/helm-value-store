package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/skuid/helm-value-store/dynamo"
	"github.com/skuid/helm-value-store/store"
	"github.com/spf13/cobra"
)

type loadCmdArgs struct {
	file        string
	table       string
	createTable bool
}

var loadArgs = &loadCmdArgs{}

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "load a json file of releases",
	Run:   load,
}

func init() {
	RootCmd.AddCommand(loadCmd)
	loadCmd.Flags().StringVar(&loadArgs.file, "file", "dynamoReleases.json", "Name of file to ingest")
	loadCmd.Flags().StringVar(&loadArgs.table, "table", "helm-charts", "Name of table")
	loadCmd.Flags().BoolVar(&loadArgs.createTable, "create-table", false, "Create the table on load")

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

	rs, err := dynamo.NewReleaseStore(loadArgs.table)
	exitOnErr(err)

	if loadArgs.createTable {
		rs.CreateTable()
	}

	err = rs.Load(releases)
	exitOnErr(err)
	fmt.Println("Loaded resources into dynamo!")
}
