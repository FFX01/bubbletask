package list

import (
	"strings"

	"github.com/FFX01/bubbletask/internal/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	selectedBackgroundColor lipgloss.Color = lipgloss.Color("3")
	selectedForegroundColor                = lipgloss.Color("4")
)

var (
	selectedStyle lipgloss.Style = lipgloss.NewStyle().
		Background(selectedBackgroundColor).
		Foreground(selectedForegroundColor)
)

type ListItem interface {
	Title() string
	Description() string
}

type Model struct {
	items       []ListItem
	selectedidx int
	isAdding    bool
	input       textinput.Model
	Title       string
	focused     bool
}

func New() Model {
	m := Model{}
	m.items = make([]ListItem, 0)
	m.input = textinput.New()
	return m
}

func (m *Model) Focus() {
	m.focused = true
}

func (m *Model) Blur() {
	m.focused = false
}

func (m *Model) SelectedIndex() int {
	return m.selectedidx
}

func (m *Model) SetItems(v ...ListItem) {
	m.items = v
}

func (m *Model) AddItem(v ListItem) {
	m.items = append(m.items, v)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) onUp() {
	if m.selectedidx == 0 {
		m.selectedidx = len(m.items) - 1
	} else {
		m.selectedidx--
	}
}

func (m *Model) onDown() {
	if m.selectedidx == len(m.items)-1 {
		m.selectedidx = 0
	} else {
		m.selectedidx++
	}
}

func (m *Model) onAddItem() {
	m.isAdding = true
	m.input.Focus()
}

type AddedItem struct {
	Value string
}

func AddedItemCmd(v string) func() tea.Msg {
	return func() tea.Msg {
		return AddedItem{Value: v}
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.onUp()
		case "down":
			m.onDown()
		case "a":
			if !m.isAdding {
				m.onAddItem()
				return m, nil
			}
		}
	case textinput.Confirm:
		m.isAdding = false
		m.input.Reset()
		m.selectedidx++
		return m, AddedItemCmd(msg.Value)
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	outputItems := []string{}

	titleStyle := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Bold(true).
		Underline(true).
		MarginBottom(1)
	outputItems = append(outputItems, titleStyle.Render(m.Title))

	for idx, item := range m.items {
		if idx == m.selectedidx && m.focused {
			outputItems = append(outputItems, selectedStyle.Render(item.Title()))
			if m.isAdding {
				outputItems = append(outputItems, m.input.View())
			}
		} else {
			outputItems = append(outputItems, item.Title())
		}
	}

	return strings.Join(outputItems, "\n")
}
