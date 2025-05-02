package catalog

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/glue"
	"github.com/frocore/fedramp-data-mesh/cli/internal/config"
	"github.com/frocore/fedramp-data-mesh/cli/internal/logging"
	"github.com/frocore/fedramp-data-mesh/cli/internal/security"
)

type Client struct {
	cfg        *config.Config
	secCtx     *security.SecurityContext
	log        *logging.Logger
	glueClient *glue.Glue
}

type DataProduct struct {
	Name         string
	Domain       string
	Description  string
	Type         string  // "event-stream", "table", etc.
	Location     string  // S3 path
	Format       string  // "iceberg", "delta", "parquet"
	Schema       string  // JSON schema representation
	Owner        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Tags         map[string]string
}

func NewClient(cfg *config.Config, secCtx *security.SecurityContext, log *logging.Logger) (*Client, error) {
	// Get AWS session
	sess, err := secCtx.GetAWSSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}
	
	// Create Glue client
	glueClient := glue.New(sess)
	
	return &Client{
		cfg:        cfg,
		secCtx:     secCtx,
		log:        log,
		glueClient: glueClient,
	}, nil
}

func (c *Client) ListDataProducts(domainFilter string) ([]string, error) {
	// Get databases (domains)
	var productNames []string
	
	// First get all databases (represents domains)
	dbInput := &glue.GetDatabasesInput{}
	err := c.glueClient.GetDatabasesPages(dbInput, func(page *glue.GetDatabasesOutput, lastPage bool) bool {
		for _, db := range page.DatabaseList {
			// Skip if domain filter is set and doesn't match
			if domainFilter != "" && *db.Name != domainFilter {
				continue
			}
			
			// For each database, get tables (data products)
			tableInput := &glue.GetTablesInput{
				DatabaseName: db.Name,
			}
			
			err := c.glueClient.GetTablesPages(tableInput, func(tablePage *glue.GetTablesOutput, tableLastPage bool) bool {
				for _, table := range tablePage.TableList {
					// Check if this is a data product by looking for specific tags
					isDataProduct := false
					if table.Parameters != nil {
						_, isDataProduct = table.Parameters["data_product"]
					}
					
					if isDataProduct {
						// Format as domain.product
						productName := fmt.Sprintf("%s.%s", *db.Name, *table.Name)
						productNames = append(productNames, productName)
					}
				}
				return true // Continue pagination
			})
			
			if err != nil {
				c.log.Errorf("Error getting tables for database %s: %v", *db.Name, err)
			}
		}
		return true // Continue pagination
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to list data products: %w", err)
	}
	
	return productNames, nil
}

func (c *Client) GetDataProduct(name string) (*DataProduct, error) {
	// Parse domain and product name
	parts := strings.Split(name, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid data product name format, expected domain.product: %s", name)
	}
	
	domain := parts[0]
	productName := parts[1]
	
	// Get table details from Glue
	tableInput := &glue.GetTableInput{
		DatabaseName: aws.String(domain),
		Name:         aws.String(productName),
	}
	
	tableOutput, err := c.glueClient.GetTable(tableInput)
	if err != nil {
		return nil, fmt.Errorf("failed to get data product details: %w", err)
	}
	
	if tableOutput.Table == nil {
		return nil, fmt.Errorf("data product not found: %s", name)
	}
	
	// Extract metadata
	table := tableOutput.Table
	
	// Check if this is a data product
	isDataProduct := false
	if table.Parameters != nil {
		_, isDataProduct = table.Parameters["data_product"]
	}
	
	if !isDataProduct {
		return nil, fmt.Errorf("table is not marked as a data product: %s", name)
	}
	
	// Build data product object
	product := &DataProduct{
		Name:        name,
		Domain:      domain,
		Description: aws.StringValue(table.Description),
		Location:    aws.StringValue(table.StorageDescriptor.Location),
		Tags:        make(map[string]string),
	}
	
	// Extract additional metadata from parameters
	if table.Parameters != nil {
		if format, ok := table.Parameters["table_format"]; ok {
			product.Format = *format
		} else {
			// Default to Iceberg if not specified
			product.Format = "iceberg"
		}
		
		if productType, ok := table.Parameters["data_product_type"]; ok {
			product.Type = *productType
		}
		
		if owner, ok := table.Parameters["owner"]; ok {
			product.Owner = *owner
		}
	}
	
	// Set timestamps
	if table.CreateTime != nil {
		product.CreatedAt = *table.CreateTime
	}
	
	if table.UpdateTime != nil {
		product.UpdatedAt = *table.UpdateTime
	}
	
	// Get tags if any
	tagsInput := &glue.GetTagsInput{
		ResourceArn: aws.String(fmt.Sprintf("arn:aws:glue:%s:%s:table/%s/%s", 
			c.cfg.AWSRegion, c.cfg.AWSAccountID, domain, productName)),
	}
	
	tagsOutput, err := c.glueClient.GetTags(tagsInput)
	if err == nil && tagsOutput.Tags != nil {
		for k, v := range tagsOutput.Tags {
			product.Tags[k] = *v
		}
	}
	
	return product, nil
}

func (c *Client) GetDataProductPath(name string) (string, error) {
	product, err := c.GetDataProduct(name)
	if err != nil {
		return "", err
	}
	
	return product.Location, nil
}

func (c *Client) GetDataProductSchema(name string) (string, error) {
	// Parse domain and product name
	parts := strings.Split(name, ".")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid data product name format, expected domain.product: %s", name)
	}
	
	domain := parts[0]
	productName := parts[1]
	
	// Get table details from Glue
	tableInput := &glue.GetTableInput{
		DatabaseName: aws.String(domain),
		Name:         aws.String(productName),
	}
	
	tableOutput, err := c.glueClient.GetTable(tableInput)
	if err != nil {
		return "", fmt.Errorf("failed to get data product schema: %w", err)
	}
	
	if tableOutput.Table == nil || tableOutput.Table.StorageDescriptor == nil {
		return "", fmt.Errorf("data product schema not found: %s", name)
	}
	
	// Convert Glue schema to JSON
	columns := tableOutput.Table.StorageDescriptor.Columns
	schema := map[string]interface{}{
		"type": "struct",
		"fields": []map[string]interface{}{},
	}
	
	fields := []map[string]interface{}{}
	for _, col := range columns {
		field := map[string]interface{}{
			"name": *col.Name,
			"type": *col.Type,
		}
		
		if col.Comment != nil {
			field["comment"] = *col.Comment
		}
		
		fields = append(fields, field)
	}
	
	schema["fields"] = fields
	
	// Convert to JSON
	schemaBytes, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal schema to JSON: %w", err)
	}
	
	return string(schemaBytes), nil
}
