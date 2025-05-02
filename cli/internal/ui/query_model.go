package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/frocore/fedramp-data-mesh/cli/internal/catalog"
	"github.com/frocore/fedramp-data-mesh/cli/internal/config"
	"github.com/frocore/fedramp-data-mesh/cli/internal/duckdb"
	"github.com/frocore/fedramp-data-mesh/cli/internal/logging"
	"github.com/frocore/fedramp-data-mesh/cli/internal/security"
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#0077B6")).
		Padding(0, 1).
		Width(80)
	
	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Padding(0, 1)
	
	promptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#0077B6")).
		Bold(true)
	
	resultHeaderStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#023E8A")).
		Bold(true).
		Padding(0, 1)
	
	resultCellStyle = lipgloss.NewStyle().
		Padding(0, 1)
)

// Model represents the UI state
type QueryModel struct {
	cfg           *config.Config
	secCtx        *security.SecurityContext
	log           *logging.Logger
	db            *duckdb.Connection
	queryInput    textinput.Model
	resultViewport viewport.Model
	dataProducts  []string
	selectedProduct string
	result        *duckdb.QueryResult
	error         string
	width         int
	height        int
}

func NewQueryModel(cfg *config.Config, secCtx *security.SecurityContext, log *logging.Logger, initialProduct string) *QueryModel {
	// Create text input for SQL queries
	ti := textinput.New()
	ti.Placeholder = "Enter SQL query"
	ti.Width = 80
	ti.Focus()
	
	// Create viewport for results
	vp := viewport.New(80, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#0077B6"))
	
	// Create DB connection
	db, err := duckdb.NewConnection(cfg, secCtx)
	if err != nil {
		log.Errorf("Failed to create DuckDB connection: %v", err)
	}
	
	// Fetch available data products
	catalogClient, err := catalog.NewClient(cfg, secCtx, log)
	if err != nil {
		log.Errorf("Failed to create catalog client: %v", err)
	}
	
	var products []string
	if catalogClient != nil {
		products, err = catalogClient.ListDataProducts("")
		if err != nil {
			log.Errorf("Failed to fetch data products: %v", err)
		}
	}
	
	selectedProduct := initialProduct
	if selectedProduct == "" && len(products) > 0 {
		selectedProduct = products[0]
	}
	
	if db != nil && catalogClient != nil && selectedProduct != "" {
		// Resolve product path
		path, err := catalogClient.GetDataProductPath(selectedProduct)
		if err != nil {
			log.Errorf("Failed to resolve data product path: %v", err)
		} else {
			// Register the data product with DuckDB
			if err := db.RegisterDataProduct(selectedProduct, path); err != nil {
				log.Errorf("Failed to register data product: %v", err)
			}
		}
	}
	
	return &QueryModel{
		cfg:            cfg,
		secCtx:         secCtx,
		log:            log,
		db:             db,
		queryInput:     ti,
		resultViewport: vp,
		dataProducts:   products,
		selectedProduct: selectedProduct,
		error:          "",
		width:          80,
		height:         24,
	}
}

// Init implements bubbletea.Model
func (m *QueryModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements bubbletea.Model
func (m *QueryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd  tea.Cmd
		vpCmd  tea.Cmd
	)
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		
		case "enter":
			if m.queryInput.Value() != "" {
				return m, m.executeQuery
			}
		
		case "tab":
			// Toggle between input and result view
			if m.queryInput.Focused() {
				m.queryInput.Blur()
				// Allow scrolling in results
			} else {
				m.queryInput.Focus()
			}
			return m, nil
		}
	
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resultViewport.Width = msg.Width
		m.resultViewport.Height = msg.Height - 10 // Leave room for query input and headers
		return m, nil
		
	case queryResultMsg:
		m.result = msg.result
		m.error = ""
		m.resultViewport.SetContent(m.formatResult())
		return m, nil
		
	case queryErrorMsg:
		m.error = msg.error
		return m, nil
	}
	
	m.queryInput, tiCmd = m.queryInput.Update(msg)
	m.resultViewport, vpCmd = m.resultViewport.Update(msg)
	
	return m, tea.Batch(tiCmd, vpCmd)
}

// View implements bubbletea.Model
func (m *QueryModel) View() string {
	var b strings.Builder
	
	// Title bar
	title := fmt.Sprintf(" FroCore Data Mesh CLI - Query Tool ")
	b.WriteString(titleStyle.Render(title))
	b.WriteString("\n\n")
	
	// Current data product
	b.WriteString(fmt.Sprintf("Current Data Product: %s\n\n", m.selectedProduct))
	
	// Query input
	b.WriteString(promptStyle.Render("SQL Query: "))
	b.WriteString("\n")
	b.WriteString(m.queryInput.View())
	b.WriteString("\n\n")
	
	// Error message (if any)
	if m.error != "" {
		b.WriteString(errorStyle.Render("Error: " + m.error))
		b.WriteString("\n\n")
	}
	
	// Results
	if m.result != nil {
		b.WriteString("Results:\n")
		b.WriteString(m.resultViewport.View())
	} else {
		b.WriteString("No results to display. Press Enter to execute query.")
	}
	
	// Help text
	b.WriteString("\n\n")
	b.WriteString("Press Tab to toggle focus, Enter to execute, Esc to quit")
	
	return b.String()
}

func (m *QueryModel) executeQuery() tea.Msg {
	if m.db == nil {
		return queryErrorMsg{error: "Database connection not initialized"}
	}
	
	query := m.queryInput.Value()
	
	// Check for data product access
	if m.selectedProduct != "" {
		canAccess, err := m.secCtx.CanAccessDataProduct(m.selectedProduct)
		if err != nil {
			return queryErrorMsg{error: fmt.Sprintf("Failed to check access: %v", err)}
		}
		if !canAccess {
			return queryErrorMsg{error: fmt.Sprintf("Access denied to data product: %s", m.selectedProduct)}
		}
	}
	
	// Execute the query
	result, err := m.db.ExecuteQuery(query)
	if err != nil {
		return queryErrorMsg{error: err.Error()}
	}
	
	return queryResultMsg{result: result}
}

func (m *QueryModel) formatResult() string {
	if m.result == nil || len(m.result.Columns) == 0 {
		return "No results to display"
	}
	
	var b strings.Builder
	
	// Calculate column widths
	colWidths := make([]int, len(m.result.Columns))
	for i, col := range m.result.Columns {
		colWidths[i] = len(col) + 2 // Add padding
	}
	
	// Check row data to ensure column width is sufficient
	for _, row := range m.result.Rows {
		for i, val := range row {
			if i < len(colWidths) {
				valStr := fmt.Sprintf("%v", val)
				if len(valStr)+2 > colWidths[i] {
					colWidths[i] = len(valStr) + 2
				}
			}
		}
	}
	
	// Header row
	for i, col := range m.result.Columns {
		b.WriteString(resultHeaderStyle.Width(colWidths[i]).Render(col))
	}
	b.WriteString("\n")
	
	// Data rows
	for _, row := range m.result.Rows {
		for i, val := range row {
			if i < len(colWidths) {
				valStr := fmt.Sprintf("%v", val)
				b.WriteString(resultCellStyle.Width(colWidths[i]).Render(valStr))
			}
		}
		b.WriteString("\n")
	}
	
	// Summary
	b.WriteString(fmt.Sprintf("\n%d rows returned", len(m.result.Rows)))
	
	return b.String()
}

// Message types for Bubble Tea
type queryResultMsg struct {
	result *duckdb.QueryResult
}

type queryErrorMsg struct {
	error string
}
