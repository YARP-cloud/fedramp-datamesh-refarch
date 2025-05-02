package cmd

import (
	"fmt"

	"github.com/frocore/fedramp-data-mesh/cli/internal/catalog"
	"github.com/frocore/fedramp-data-mesh/cli/internal/config"
	"github.com/frocore/fedramp-data-mesh/cli/internal/logging"
	"github.com/frocore/fedramp-data-mesh/cli/internal/security"
	"github.com/spf13/cobra"
)

func NewInfoCmd(cfg *config.Config, secCtx *security.SecurityContext, log *logging.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info [data_product]",
		Short: "Show information about a data product",
		Long:  `Display detailed information about a specific data product`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dataProduct := args[0]
			return getInfo(dataProduct, cfg, secCtx, log)
		},
	}
	
	return cmd
}

func getInfo(dataProduct string, cfg *config.Config, secCtx *security.SecurityContext, log *logging.Logger) error {
	catalogClient, err := catalog.NewClient(cfg, secCtx, log)
	if err != nil {
		return err
	}
	
	product, err := catalogClient.GetDataProduct(dataProduct)
	if err != nil {
		return err
	}
	
	fmt.Println("Data Product Information")
	fmt.Println("=======================")
	fmt.Printf("Name:         %s\n", product.Name)
	fmt.Printf("Domain:       %s\n", product.Domain)
	fmt.Printf("Description:  %s\n", product.Description)
	fmt.Printf("Type:         %s\n", product.Type)
	fmt.Printf("Format:       %s\n", product.Format)
	fmt.Printf("Location:     %s\n", product.Location)
	fmt.Printf("Owner:        %s\n", product.Owner)
	fmt.Printf("Created:      %s\n", product.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Last Updated: %s\n", product.UpdatedAt.Format("2006-01-02 15:04:05"))
	
	fmt.Println("\nTags:")
	for k, v := range product.Tags {
		fmt.Printf("  %s: %s\n", k, v)
	}
	
	return nil
}
