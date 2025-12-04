package finder

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type selectorModel struct {
	allItems      []entry
	filteredItems []entry
	cursor        int
	input         textinput.Model
	width         int
	height        int
	selected      *entry
	cancelled     bool
}

func newSelector(items []entry) selectorModel {
	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.Prompt = "> "
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 50

	return selectorModel{
		allItems:      items,
		filteredItems: items,
		input:         ti,
		cursor:        0,
	}
}

func (m selectorModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m selectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.cancelled = true
			return m, tea.Quit

		case "enter":
			if len(m.filteredItems) > 0 {
				m.selected = &m.filteredItems[m.cursor]
			}
			return m, tea.Quit

		case "up", "ctrl+p":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "ctrl+n":
			if m.cursor < len(m.filteredItems)-1 {
				m.cursor++
			}

		default:
			m.input, cmd = m.input.Update(msg)
			m.filterItems()
			// Reset cursor if out of bounds
			if m.cursor >= len(m.filteredItems) {
				m.cursor = 0
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, cmd
}

func (m *selectorModel) filterItems() {
	query := strings.TrimSpace(m.input.Value())
	if query == "" {
		m.filteredItems = m.allItems
		return
	}

	// Split by space for AND search
	keywords := strings.Fields(strings.ToLower(query))
	if len(keywords) == 0 {
		m.filteredItems = m.allItems
		return
	}

	var filtered []entry
	for _, item := range m.allItems {
		label := strings.ToLower(item.displayLabel)
		match := true
		for _, kw := range keywords {
			if !strings.Contains(label, kw) {
				match = false
				break
			}
		}
		if match {
			filtered = append(filtered, item)
		}
	}

	m.filteredItems = filtered
}

func (m selectorModel) View() string {
	var b strings.Builder

	// Input field at top
	b.WriteString(m.input.View() + "\n")

	// Show count
	countStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))
	b.WriteString("  " + countStyle.Render(fmt.Sprintf("%d/%d", len(m.filteredItems), len(m.allItems))) + "\n")

	// List items
	visibleHeight := m.height - 5 // Reserve space for input and count
	if visibleHeight < 1 {
		visibleHeight = 10
	}

	start := m.cursor - visibleHeight/2
	if start < 0 {
		start = 0
	}
	end := start + visibleHeight
	if end > len(m.filteredItems) {
		end = len(m.filteredItems)
		start = end - visibleHeight
		if start < 0 {
			start = 0
		}
	}

	selectedStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230"))

	normalStyle := lipgloss.NewStyle()

	for i := start; i < end; i++ {
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}

		line := cursor + m.filteredItems[i].displayLabel

		if i == m.cursor {
			b.WriteString(selectedStyle.Render(line) + "\n")
		} else {
			b.WriteString(normalStyle.Render(line) + "\n")
		}
	}

	if len(m.filteredItems) == 0 {
		noResultStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Italic(true)
		b.WriteString("  " + noResultStyle.Render("No matches found\n"))
	}

	return b.String()
}

func runSelector(items []entry) (*entry, error) {
	p := tea.NewProgram(newSelector(items), tea.WithAltScreen())

	model, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("error running selector: %w", err)
	}

	m := model.(selectorModel)
	if m.cancelled {
		return nil, ErrCancelled
	}

	return m.selected, nil
}
