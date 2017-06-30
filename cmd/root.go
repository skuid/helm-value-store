package cmd

import (
	"fmt"
	"os"

	"github.com/skuid/helm-value-store/dynamo"
	"github.com/skuid/helm-value-store/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var releaseStore store.ReleaseStore

var RootCmd = &cobra.Command{
	Use:   "helm-value-store",
	Short: "A helm plugin for working with Helm Release data",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		switch backend := viper.GetString("backend"); backend {
		case "dynamodb":
			releaseStore, err = dynamo.NewReleaseStore(viper.GetString("dynamodb-table"))
		default:
			err = fmt.Errorf("No valid value store specified! '%s'", backend)
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

	RootCmd.PersistentFlags().String("backend", "dynamodb", "The backend for the value store")

	// DynamoDB flags
	RootCmd.PersistentFlags().String("dynamodb-table", "helm-charts", "Name of the dynamodb table")
}

func initConfig() {
	if err := viper.BindPFlags(RootCmd.PersistentFlags()); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	viper.SetEnvPrefix("HELM_VALUE_STORE")
	viper.AutomaticEnv()
}

var valueExtensions = []string{"json", "yaml", "yml"}
