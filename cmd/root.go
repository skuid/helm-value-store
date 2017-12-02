package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/skuid/helm-value-store/datastore"
	"github.com/skuid/helm-value-store/dynamo"
	"github.com/skuid/helm-value-store/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var releaseStore store.ReleaseStore

var storeTypes = []string{"dynamodb", "datastore"}

// RootCmd is the root command
var RootCmd = &cobra.Command{
	Use:   "helm-value-store",
	Short: "A helm plugin for working with Helm Release data",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		switch backend := viper.GetString("backend"); backend {
		case "dynamodb":
			releaseStore, err = dynamo.NewReleaseStore(viper.GetString("dynamodb-table"))
		case "datastore":
			releaseStore, err = datastore.NewReleaseStore(viper.GetString("service-account"))
		default:
			err = fmt.Errorf("No valid value store specified: %s. Must be one of %v", backend, storeTypes)
		}
		exitOnErr(err)
	},
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func init() {
	err := os.Setenv("AWS_SDK_LOAD_CONFIG", "true")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().String("backend", "dynamodb", fmt.Sprintf("The backend for the value store. Must be one of %v", storeTypes))

	// DynamoDB flags
	RootCmd.PersistentFlags().String("dynamodb-table", "helm-charts", "Name of the dynamodb table")
	RootCmd.PersistentFlags().String("service-account", "sa.json", "The Google Service Account JSON file")
	RootCmd.PersistentFlags().Duration("timeout", time.Duration(30)*time.Second, "The timeout for a given command")
}

func initConfig() {
	if err := viper.BindPFlags(RootCmd.PersistentFlags()); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	viper.BindPFlags(RootCmd.Flags())
	viper.SetEnvPrefix("HELM_VALUE_STORE")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

var valueExtensions = []string{"json", "yaml", "yml"}
