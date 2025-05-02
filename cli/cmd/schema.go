package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/frocore/fedramp-data-mesh/cli/internal/catalog"
	"github.com/frocore/fedramp-data-mesh/cli/internal/config"
	"github.com/frocore/fedramp-data-mesh/cli/internal/logging"
	"github.com/frocore/fedramp-data-mesh/cli/internal/security"
	"github.com/spf13/cobra"
)

func NewSchemaCmd(cfg *config.Config, secCtx *security.SecurityContext, log *logging.Logger) *cobra.Command {
	var outputFile string
	var formatOutput bool
	
	cmd := &cobra.Command{
		Use:   "schema [data_product]",
		Short: "Show schema for a data product",
		Long:  `Display the schema definition for a specific data product`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dataProduct := args[0]
			return getSchema(dataProduct, outputFile, formatOutput, cfg, secCtx, log)
		},
	}
	
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write schema to file")
	cmd.Flags().BoolVarP(&formatOutput, "format", "f", true, "Format JSON output")
	
	return cmd
}

func getSchema(dataProduct, outputFile string, formatOutput bool, cfg *config.Config, secCtx *security.SecurityContext, log *logging.Logger) error {
	catalogClient, err := catalog.NewClient(cfg, secCtx, log)
	if err != nil {
		return err
	}
	
	schema, err := catalogClient.GetDataProductSchema(dataProduct)
	if err != nil {
		return err
	}
	
	var output []byte
	if formatOutput {
		var jsonObj interface{}
		if err := json.Unmarshal([]byte(schema), &jsonObj); err != nil {
			return err
		}
		output, err = json.MarshalIndent(jsonObj, "", "  ")
		if err != nil {
			return err
		}
	} else {
		output = []byte(schema)
	}
	
	if outputFile != "" {
		return os.WriteFile(outputFile, output, 0644)
	}
	
	fmt.Println(string(output))
	return nil
}
