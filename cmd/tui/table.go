package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ppmpreetham/vesper/tools"
)

// TableModel represents the results table
type TableModel struct {
	table     table.Model
	width     int
	height    int
	rowOffset int
}

func NewTableModel(results []tools.ReturnData) TableModel {
	// Create columns
	columns := []table.Column{
		{Title: "Site", Width: 30},
		{Title: "URL", Width: 70},
	}

	// Create rows from results
	rows := []table.Row{}
	for _, result := range results {
		rows = append(rows, table.Row{result.Name, result.URL})
	}

	// Create table
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(20), // Show more rows at a time
	)

	// Style the table
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return TableModel{table: t}
}

func (m TableModel) Init() tea.Cmd {
	return nil
}

func (m TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Adjust table height to fit window
		headerHeight := 4
		footerHeight := 2
		maxTableHeight := m.height - headerHeight - footerHeight

		if maxTableHeight > 0 && maxTableHeight < len(m.table.Rows()) {
			m.table.SetHeight(maxTableHeight)
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m TableModel) View() string {
	baseStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1, 2)

	title := titleStyle.Render(" FOUND PROFILES ")

	countStr := fmt.Sprintf("Found %d profiles", len(m.table.Rows()))
	help := "\n↑/↓: Navigate • Page Up/Down: Scroll • q: Quit\n"

	return title + "\n" + countStr + "\n" + baseStyle.Render(m.table.View()) + help
}
