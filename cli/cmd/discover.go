package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/frocore/fedramp-data-mesh/cli/internal/catalog"
	"github.com/frocore/fedramp-data-mesh/cli/internal/config"
	"github.com/frocore/fedramp-data-mesh/cli/internal/logging"
	"github.com/frocore/fedramp-data-mesh/cli/internal/security"
	"github.com/frocore/fedramp-data-mesh/cli/internal/ui"
	"github.com/spf13/cobra"
)

func NewDiscoverCmd(cfg *config.Config, secCtx *security.SecurityContext, log *logging.Logger) *cobra.Command {
	var domainFilter string
	var interactive bool

	cmd := &cobra.Command{
		Use:   "discover",
		Short: "Discover available data products",
		Long:  `Browse and discover available data products in the data mesh`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if interactive {
				return launchDiscoverUI(cfg, secCtx, log, domainFilter)
			} else {
				return listDataProducts(cfg, secCtx, log, domainFilter)
			}
		},
	}

	cmd.Flags().StringVarP(&domainFilter, "domain", "d", "", "Filter by domain")
	cmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Use interactive mode")

	return cmd
}

func listDataProducts(cfg *config.Config, secCtx *security.SecurityContext, log *logging.Logger, domainFilter string) error {
	catalogClient, err := catalog.NewClient(cfg, secCtx, log)
	if err != nil {
		return err
	}

	products, err := catalogClient.ListDataProducts(domainFilter)
	if err != nil {
		return err
	}

	if len(products) == 0 {
		fmt.Println("No data products found")
		return nil
	}

	fmt.Println("Available data products:")
	fmt.Println("")

	for _, product := range products {
		fmt.Printf("- %s\n", product)
	}

	return nil
}

func launchDiscoverUI(cfg *config.Config, secCtx *security.SecurityContext, log *logging.Logger, domainFilter string) error {
	model := ui.NewDiscoverModel(cfg, secCtx, log, domainFilter)
	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
