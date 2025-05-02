package cmd

import (
	"github.com/frocore/fedramp-data-mesh/cli/internal/config"
	"github.com/frocore/fedramp-data-mesh/cli/internal/logging"
	"github.com/frocore/fedramp-data-mesh/cli/internal/security"
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

func Execute(cfg *config.Config, secCtx *security.SecurityContext, log *logging.Logger) error {
	rootCmd = &cobra.Command{
		Use:   "dmesh",
		Short: "FroCore Data Mesh CLI",
		Long:  `Command-line tool for interacting with the FroCore Event-Driven Data Mesh`,
	}
	
	// Add subcommands
	rootCmd.AddCommand(NewDiscoverCmd(cfg, secCtx, log))
	rootCmd.AddCommand(NewQueryCmd(cfg, secCtx, log))
	rootCmd.AddCommand(NewSchemaCmd(cfg, secCtx, log))
	rootCmd.AddCommand(NewInfoCmd(cfg, secCtx, log))
	
	return rootCmd.Execute()
}
