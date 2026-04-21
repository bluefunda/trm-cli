package cmd

import (
	"github.com/spf13/cobra"

	"github.com/bluefunda/trm-cli/internal/config"
	"github.com/bluefunda/trm-cli/internal/ui"
)

var (
	cfgBFF    string
	cfgDomain string
	cfgOutput string
)

var Version = "dev"

var rootCmd = &cobra.Command{
	Use:     "trm",
	Short:   "TRM -- CLI for the bluerequests change/release management platform",
	Long:    "TRM is a command-line interface for interacting with the TRM platform via gRPC.",
	Version: Version,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgBFF, "bff", "", "BFF gRPC address host:port (overrides config)")
	rootCmd.PersistentFlags().StringVar(&cfgDomain, "domain", "", "Domain (overrides config)")
	rootCmd.PersistentFlags().StringVarP(&cfgOutput, "output", "o", "", "Output format: table, json, quiet")

	rootCmd.AddCommand(
		loginCmd,
		healthCmd,
		versionCmd,
		userCmd,
		eventsCmd,
		rpcCmd,
	)
}

func Execute() error {
	return rootCmd.Execute()
}

func loadConfig() *config.Config {
	cfg, err := config.Load()
	if err != nil {
		ui.Error("Failed to load config: " + err.Error())
		return &config.Config{}
	}
	if cfgBFF != "" {
		cfg.BFFURL = cfgBFF
	}
	if cfgDomain != "" {
		cfg.Domain = cfgDomain
	}
	return cfg
}

func outputFormat(cfg *config.Config) ui.OutputFormat {
	if cfgOutput != "" {
		switch cfgOutput {
		case "json":
			return ui.FormatJSON
		case "quiet":
			return ui.FormatQuiet
		default:
			return ui.FormatTable
		}
	}
	switch cfg.Defaults.Output {
	case "json":
		return ui.FormatJSON
	case "quiet":
		return ui.FormatQuiet
	default:
		return ui.FormatTable
	}
}
