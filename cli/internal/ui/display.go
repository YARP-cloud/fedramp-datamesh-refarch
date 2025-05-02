package ui

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/frocore/fedramp-data-mesh/cli/internal/duckdb"
)

func DisplayQueryResults(result *duckdb.QueryResult, format string) error {
	switch strings.ToLower(format) {
	case "table":
		return displayTable(result)
	case "csv":
		return displayCSV(result)
	case "json":
		return displayJSON(result)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

func displayTable(result *duckdb.QueryResult) error {
	if result == nil || len(result.Columns) == 0 {
		fmt.Println("No results to display")
		return nil
	}
	
	// Calculate column widths
	colWidths := make([]int, len(result.Columns))
	for i, col := range result.Columns {
		colWidths[i] = len(col)
	}
	
	// Check row data to ensure column width is sufficient
	for _, row := range result.Rows {
		for i, val := range row {
			if i < len(colWidths) {
				valStr := fmt.Sprintf("%v", val)
				if len(valStr) > colWidths[i] {
					colWidths[i] = len(valStr)
				}
			}
		}
	}
	
	// Print header
	for i, col := range result.Columns {
		fmt.Printf("| %-*s ", colWidths[i], col)
	}
	fmt.Println("|")
	
	// Print separator
	for i := range result.Columns {
		fmt.Print("|")
		for j := 0; j < colWidths[i]+2; j++ {
			fmt.Print("-")
		}
	}
	fmt.Println("|")
	
	// Print data rows
	for _, row := range result.Rows {
		for i, val := range row {
			if i < len(colWidths) {
				valStr := fmt.Sprintf("%v", val)
				fmt.Printf("| %-*s ", colWidths[i], valStr)
			}
		}
		fmt.Println("|")
	}
	
	// Print footer
	fmt.Printf("\n%d rows returned\n", len(result.Rows))
	
	return nil
}

func displayCSV(result *duckdb.QueryResult) error {
	if result == nil || len(result.Columns) == 0 {
		fmt.Println("No results to display")
		return nil
	}
	
	w := csv.NewWriter(os.Stdout)
	
	// Write header
	if err := w.Write(result.Columns); err != nil {
		return fmt.Errorf("error writing CSV header: %w", err)
	}
	
	// Write data rows
	for _, row := range result.Rows {
		// Convert row to []string
		rowStrings := make([]string, len(row))
		for i, val := range row {
			rowStrings[i] = fmt.Sprintf("%v", val)
		}
		
		if err := w.Write(rowStrings); err != nil {
			return fmt.Errorf("error writing CSV row: %w", err)
		}
	}
	
	w.Flush()
	
	if err := w.Error(); err != nil {
		return fmt.Errorf("error flushing CSV writer: %w", err)
	}
	
	return nil
}

func displayJSON(result *duckdb.QueryResult) error {
	if result == nil || len(result.Columns) == 0 {
		fmt.Println("No results to display")
		return nil
	}
	
	// Convert to list of maps
	data := make([]map[string]interface{}, len(result.Rows))
	
	for i, row := range result.Rows {
		item := make(map[string]interface{})
		
		for j, val := range row {
			if j < len(result.Columns) {
				item[result.Columns[j]] = val
			}
		}
		
		data[i] = item
	}
	
	// Marshal to JSON
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling to JSON: %w", err)
	}
	
	fmt.Println(string(jsonBytes))
	
	return nil
}
