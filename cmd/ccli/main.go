package main

import (
	"fmt"
	"go/doc"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"gitlab.corp.mail.ru/cloud/go/cli/internal/cmd/completion"
	"gitlab.corp.mail.ru/cloud/go/cli/internal/scli"
)

func main() {
	cmd := &cobra.Command{
		Use:           "scli",
		Short:         "Cloud command line interface",
		Long:          "Cloud command line interface",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	configFile := cmd.PersistentFlags().StringP("config", "c", "", "Config file (default is ~/.scli, then /etc/scli/config.yaml)")
	verbose := cmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose mode")

	cmd.PersistentFlags().StringP("environment", "e", "", "Current environment (prod,dev,...) (required)")
	viper.BindPFlag("environment", cmd.PersistentFlags().Lookup("environment"))

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if !*verbose {
			log.SetOutput(ioutil.Discard)
		}

		cfg, err := scli.GetConfigPath(*configFile)
		if err != nil {
			return err
		}

		log.Infof("Start using config %s", cfg)

		viper.SetConfigType("yaml")
		viper.SetConfigFile(cfg)

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err != nil {
			return err
		}

		environment := viper.GetString("environment")
		if environment == "" {
			return fmt.Errorf("set environment in -e flag or [environment] config file")
		}

		if viper.Get(environment) == nil {
			return fmt.Errorf("environment %s  not found in config", environment)
		}
		log.Infof("Start using environment %s", environment)

		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknown-host"
		}
		hostname = strings.Split(hostname, ".")[0]

		// A bit of dark magic
		// Otherwise do not unmarshal nested structures after call viper.Set
		viper.MergeConfigMap(map[string]interface{}{
			environment: map[string]interface{}{
				"source": fmt.Sprintf("scli@%s", hostname),
			},
		})
		// End of dark magic

		return nil
	}

	cmd.AddCommand(completion.New())
	cmd.AddCommand(doc.New())
	cmd.AddCommand(user.New(&user.Env{Printer: &scli.Printer{}}))
	cmd.AddCommand(pairdb.New(&pairdb.Env{Printer: &scli.Printer{}}))
	cmd.AddCommand(filedb.New(&filedb.Env{Printer: &scli.Printer{}}))
	cmd.AddCommand(cloud.New(&cloud.Env{Printer: &scli.Printer{}}))
	cmd.AddCommand(billing.New(&billing.Env{}))
	cmd.AddCommand(etcd.New(&etcd.Env{}))
	cmd.AddCommand(weblink.New(&weblink.Env{Printer: &scli.Printer{}}))
	cmd.AddCommand(swa.New())

	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, scli.Color.Red("E: %+v\n"), err)
		os.Exit(1)
	}
}
