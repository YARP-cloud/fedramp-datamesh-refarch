package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/frocore/fedramp-data-mesh/cli/internal/catalog"
	"github.com/frocore/fedramp-data-mesh/cli/internal/config"
	"github.com/frocore/fedramp-data-mesh/cli/internal/logging"
	"github.com/frocore/fedramp-data-mesh/cli/internal/security"
)

// Item represents a data product in the list
type Item struct {
	name  string
	desc  string
}

func (i Item) Title() string       { return i.name }
func (i Item) Description() string { return i.desc }
func (i Item) FilterValue() string { return i.name }

// Model represents the UI state
type DiscoverModel struct {
	cfg           *config.Config
	secCtx        *security.SecurityContext
	log           *logging.Logger
	list          list.Model
	domainFilter  string
	selectedProduct *catalog.DataProduct
	detailViewport viewport.Model
	showingDetails bool
	err           error
	width         int
	height        int
}

func NewDiscoverModel(cfg *config.Config, secCtx *security.SecurityContext, log *logging.Logger, domainFilter string) *DiscoverModel {
	// Create list
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Available Data Products"
	l.SetShowHelp(true)
	
	// Create viewport for details
	vp := viewport.New(80, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#0077B6"))
	
	return &DiscoverModel{
		cfg:           cfg,
		secCtx:        secCtx,
		log:           log,
		list:          l,
		domainFilter:  domainFilter,
		detailViewport: vp,
		showingDetails: false,
		width:         80,
		height:        24,
	}
}

// Init implements bubbletea.Model
func (m *DiscoverModel) Init() tea.Cmd {
	return m.loadDataProducts
}

// Update implements bubbletea.Model
func (m *DiscoverModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		if m.showingDetails {
			m.detailViewport.Width = msg.Width
			m.detailViewport.Height = msg.Height
			vp, cmd := m.detailViewport.Update(msg)
			m.detailViewport = vp
			cmds = append(cmds, cmd)
		} else {
			m.list.SetSize(msg.Width, msg.Height)
			l, cmd := m.list.Update(msg)
			m.list = l
			cmds = append(cmds, cmd)
		}
		
	case tea.KeyMsg:
		if m.showingDetails {
			switch {
			case key.Matches(msg, key.NewBinding(key.WithKeys("esc", "backspace"))):
				m.showingDetails = false
				m.selectedProduct = nil
				return m, nil
			default:
				vp, cmd := m.detailViewport.Update(msg)
				m.detailViewport = vp
				cmds = append(cmds, cmd)
			}
		} else {
			switch {
			case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
				if i, ok := m.list.SelectedItem().(Item); ok {
					return m, m.loadDataProductDetails(i.name)
				}
			default:
				l, cmd := m.list.Update(msg)
				m.list = l
				cmds = append(cmds, cmd)
			}
		}
		
	case dataProductsLoadedMsg:
		items := make([]list.Item, len(msg.products))
		for i, product := range msg.products {
			items[i] = Item{name: product, desc: fmt.Sprintf("Data product: %s", product)}
		}
		m.list.SetItems(items)
		l, cmd := m.list.Update(nil)
		m.list = l
		cmds = append(cmds, cmd)
		
	case dataProductDetailsLoadedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.selectedProduct = msg.product
		m.showingDetails = true
		m.detailViewport.SetContent(m.formatProductDetails())
		vp, cmd := m.detailViewport.Update(nil)
		m.detailViewport = vp
		cmds = append(cmds, cmd)
		
	case errMsg:
		m.err = msg.err
		return m, nil
	}
	
	return m, tea.Batch(cmds...)
}

// View implements bubbletea.Model
func (m *DiscoverModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress any key to exit", m.err)
	}
	
	if m.showingDetails && m.selectedProduct != nil {
		return fmt.Sprintf("Data Product Details: %s\n\n%s\n\nPress ESC to go back", 
			m.selectedProduct.Name, m.detailViewport.View())
	}
	
	return m.list.View()
}

func (m *DiscoverModel) loadDataProducts() tea.Msg {
	catalogClient, err := catalog.NewClient(m.cfg, m.secCtx, m.log)
	if err != nil {
		return errMsg{err: fmt.Errorf("failed to create catalog client: %w", err)}
	}
	
	products, err := catalogClient.ListDataProducts(m.domainFilter)
	if err != nil {
		return errMsg{err: fmt.Errorf("failed to list data products: %w", err)}
	}
	
	return dataProductsLoadedMsg{products: products}
}

func (m *DiscoverModel) loadDataProductDetails(name string) tea.Cmd {
	return func() tea.Msg {
		catalogClient, err := catalog.NewClient(m.cfg, m.secCtx, m.log)
		if err != nil {
			return errMsg{err: fmt.Errorf("failed to create catalog client: %w", err)}
		}
		
		product, err := catalogClient.GetDataProduct(name)
		if err != nil {
			return dataProductDetailsLoadedMsg{err: err}
		}
		
		return dataProductDetailsLoadedMsg{product: product}
	}
}

func (m *DiscoverModel) formatProductDetails() string {
	if m.selectedProduct == nil {
		return "No product selected"
	}
	
	var b strings.Builder
	
	b.WriteString(fmt.Sprintf("Name:         %s\n", m.selectedProduct.Name))
	b.WriteString(fmt.Sprintf("Domain:       %s\n", m.selectedProduct.Domain))
	b.WriteString(fmt.Sprintf("Description:  %s\n", m.selectedProduct.Description))
	b.WriteString(fmt.Sprintf("Type:         %s\n", m.selectedProduct.Type))
	b.WriteString(fmt.Sprintf("Format:       %s\n", m.selectedProduct.Format))
	b.WriteString(fmt.Sprintf("Location:     %s\n", m.selectedProduct.Location))
	b.WriteString(fmt.Sprintf("Owner:        %s\n", m.selectedProduct.Owner))
	b.WriteString(fmt.Sprintf("Created:      %s\n", m.selectedProduct.CreatedAt.Format("2006-01-02 15:04:05")))
	b.WriteString(fmt.Sprintf("Last Updated: %s\n", m.selectedProduct.UpdatedAt.Format("2006-01-02 15:04:05")))
	
	b.WriteString("\nTags:\n")
	for k, v := range m.selectedProduct.Tags {
		b.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
	}
	
	return b.String()
}

// Message types for Bubble Tea
type dataProductsLoadedMsg struct {
	products []string
}

type dataProductDetailsLoadedMsg struct {
	product *catalog.DataProduct
	err     error
}

type errMsg struct {
	err error
}
