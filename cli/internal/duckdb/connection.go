package duckdb

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/marcboeker/go-duckdb"
	"github.com/frocore/fedramp-data-mesh/cli/internal/config"
	"github.com/frocore/fedramp-data-mesh/cli/internal/security"
)

type Connection struct {
	db     *sql.DB
	secCtx *security.SecurityContext
	cfg    *config.Config
}

type QueryResult struct {
	Columns []string
	Rows    [][]interface{}
}

func NewConnection(cfg *config.Config, secCtx *security.SecurityContext) (*Connection, error) {
	// Create a new in-memory DuckDB connection
	db, err := sql.Open("duckdb", "")
	if err != nil {
		return nil, fmt.Errorf("failed to open DuckDB connection: %w", err)
	}
	
	// Configure AWS credentials for S3 access
	awsAccessKey, awsSecretKey, awsSessionToken, err := secCtx.GetAWSCredentials()
	if err != nil {
		return nil, fmt.Errorf("failed to get AWS credentials: %w", err)
	}
	
	// Set up AWS credentials in DuckDB
	_, err = db.Exec(fmt.Sprintf(`
		INSTALL httpfs;
		LOAD httpfs;
		SET s3_region='%s';
		SET s3_access_key_id='%s';
		SET s3_secret_access_key='%s';
		SET s3_session_token='%s';
	`, cfg.AWSRegion, awsAccessKey, awsSecretKey, awsSessionToken))
	
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to configure S3 access in DuckDB: %w", err)
	}
	
	return &Connection{
		db:     db,
		secCtx: secCtx,
		cfg:    cfg,
	}, nil
}

func (c *Connection) Close() error {
	return c.db.Close()
}

func (c *Connection) RegisterDataProduct(name, path string) error {
	// Determine if this is Iceberg, Delta, or plain Parquet
	if strings.Contains(path, "iceberg") {
		// Iceberg table
		_, err := c.db.Exec(fmt.Sprintf(`
			INSTALL iceberg;
			LOAD iceberg;
			CREATE VIEW %s AS SELECT * FROM iceberg_scan('%s');
		`, name, path))
		return err
	} else if strings.Contains(path, "delta") {
		// Delta Lake table
		_, err := c.db.Exec(fmt.Sprintf(`
			INSTALL delta;
			LOAD delta;
			CREATE VIEW %s AS SELECT * FROM delta_scan('%s');
		`, name, path))
		return err
	} else {
		// Plain Parquet files
		_, err := c.db.Exec(fmt.Sprintf(`
			CREATE VIEW %s AS SELECT * FROM parquet_scan('%s/*.parquet');
		`, name, path))
		return err
	}
}

func (c *Connection) ExecuteQuery(query string) (*QueryResult, error) {
	// Execute the query
	rows, err := c.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()
	
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get column info: %w", err)
	}
	
	// Prepare result container
	result := &QueryResult{
		Columns: columns,
		Rows:    make([][]interface{}, 0),
	}
	
	// Prepare value containers
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	
	// Fetch rows
	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		
		rowCopy := make([]interface{}, len(values))
		for i, v := range values {
			rowCopy[i] = v
		}
		
		result.Rows = append(result.Rows, rowCopy)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}
	
	return result, nil
}
