package cmd

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/skuid/go-middlewares/authn/google"
	"github.com/skuid/helm-value-store/server"
	"github.com/skuid/spec"
	"github.com/skuid/spec/lifecycle"
	"github.com/skuid/spec/middlewares"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var level = zapcore.InfoLevel

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		l, err := spec.NewStandardLevelLogger(level)
		if err != nil {
			return fmt.Errorf("Error initializing logger: %q", err)
		}
		zap.ReplaceGlobals(l)
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		middlewareList := []middlewares.Middleware{middlewares.InstrumentRoute()}
		loggingClosures := []func(*http.Request) []zapcore.Field{}
		serverOpts := []server.ControllerOpt{}

		if viper.GetBool("auth-enabled") {
			authorizer := google.New(google.WithAuthorizedDomains(viper.GetString("email-domain")))
			serverOpts = append(serverOpts, server.WithAuthorizers(authorizer))
			middlewareList = append([]middlewares.Middleware{authorizer.Authorize()}, middlewareList...)
			loggingClosures = append(loggingClosures, authorizer.LoggingClosure)
		}

		apiController := server.NewApiController(releaseStore, serverOpts...)
		middlewareList = append(middlewareList, middlewares.Logging(loggingClosures...))

		authMux := http.NewServeMux()
		authMux.HandleFunc("/apply", apiController.ApplyChart)

		mux := http.NewServeMux()
		mux.Handle("/", middlewares.Apply(authMux, middlewareList...))

		go spec.MetricsServer(viper.GetInt("metrics-port"))

		hostPort := fmt.Sprintf(":%d", viper.GetInt("port"))
		httpServer := &http.Server{Addr: hostPort, Handler: mux}
		lifecycle.ShutdownOnTerm(httpServer)

		zap.L().Info("Starting helm-value-store server 🃏 ", zap.Int("port", viper.GetInt("port")))
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			zap.L().Fatal("Error listening", zap.Error(err))
		}
		zap.L().Info("Server gracefully stopped")

	},
}

func init() {
	RootCmd.AddCommand(serveCmd)

	// Hack to make level work
	set := flag.NewFlagSet("temp", flag.ExitOnError)
	set.Var(&level, "level", "Log level")
	levelPFlag := pflag.PFlagFromGoFlag(set.Lookup("level"))
	levelPFlag.Shorthand = "l"

	localFlagSet := serveCmd.Flags()
	localFlagSet.AddFlag(levelPFlag)
	localFlagSet.Int64P("timeout", "t", 300, "Time in seconds to timeout on installation/update of releases")
	localFlagSet.IntP("port", "p", 3000, "The port to listen on")
	localFlagSet.Int("metrics-port", 3001, "The port to listen on for metrics/health checks")
	localFlagSet.String("email-domain", "", "The email domain to filter on")
	localFlagSet.Bool("auth-enabled", true, "Enable authentication/authorization")

	viper.BindPFlags(localFlagSet)
}
