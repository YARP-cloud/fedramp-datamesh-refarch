package cmd

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/frocore/fedramp-data-mesh/cli/internal/catalog"
	"github.com/frocore/fedramp-data-mesh/cli/internal/config"
	"github.com/frocore/fedramp-data-mesh/cli/internal/duckdb"
	"github.com/frocore/fedramp-data-mesh/cli/internal/logging"
	"github.com/frocore/fedramp-data-mesh/cli/internal/security"
	"github.com/frocore/fedramp-data-mesh/cli/internal/ui"
	"github.com/spf13/cobra"
)

func NewQueryCmd(cfg *config.Config, secCtx *security.SecurityContext, log *logging.Logger) *cobra.Command {
	var dataProduct string
	var outputFormat string
	
	cmd := &cobra.Command{
		Use:   "query",
		Short: "Query data products using DuckDB",
		Long:  `Execute SQL queries against data products using DuckDB`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				// Direct query from command line
				query := args[0]
				return executeQuery(query, dataProduct, outputFormat, cfg, secCtx, log)
			} else {
				// Launch interactive UI
				return launchQueryUI(dataProduct, cfg, secCtx, log)
			}
		},
	}
	
	cmd.Flags().StringVarP(&dataProduct, "product", "p", "", "Data product to query")
	cmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table, csv, json)")
	
	return cmd
}

func executeQuery(query, dataProduct, outputFormat string, cfg *config.Config, secCtx *security.SecurityContext, log *logging.Logger) error {
	// Initialize DuckDB connection
	db, err := duckdb.NewConnection(cfg, secCtx)
	if err != nil {
		return err
	}
	defer db.Close()
	
	// Resolve full path for data product
	productPath, err := resolveDataProductPath(dataProduct, cfg, secCtx, log)
	if err != nil {
		return err
	}
	
	// Register data product in DuckDB
	if err := db.RegisterDataProduct(dataProduct, productPath); err != nil {
		return err
	}
	
	// Execute query
	result, err := db.ExecuteQuery(query)
	if err != nil {
		return err
	}
	
	// Format and display results
	return ui.DisplayQueryResults(result, outputFormat)
}

func launchQueryUI(dataProduct string, cfg *config.Config, secCtx *security.SecurityContext, log *logging.Logger) error {
	// Initialize model for Bubble Tea UI
	model := ui.NewQueryModel(cfg, secCtx, log, dataProduct)
	
	// Start the UI
	p := tea.NewProgram(model, tea.WithAltScreen())
	return p.Start()
}

func resolveDataProductPath(dataProduct string, cfg *config.Config, secCtx *security.SecurityContext, log *logging.Logger) (string, error) {
	// Query the data catalog to get the S3 path for the data product
	catalogClient, err := catalog.NewClient(cfg, secCtx, log)
	if err != nil {
		return "", err
	}
	
	return catalogClient.GetDataProductPath(dataProduct)
}
