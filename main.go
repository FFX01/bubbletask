package main

import (
	"github.com/FFX01/bubbletask/internal/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Status int

const (
	todo Status = iota
	inProgress
	done
)

var (
	focusedBorderColor lipgloss.Color = lipgloss.Color("6")
	normalBorderColor                 = lipgloss.Color("7")
)

var (
	baseListStyle lipgloss.Style = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true).
			BorderForeground(normalBorderColor).
			Padding(0, 1).
			Margin(1)
	focusedListStyle = baseListStyle.BorderForeground(focusedBorderColor)
)

type TodoItem struct {
	title       string
	description string
	status      Status
}

func (self TodoItem) Title() string       { return self.title }
func (self TodoItem) Description() string { return self.description }

type Model struct {
	todoItems      map[Status][]TodoItem
	todoList       list.Model
	inProgressList list.Model
	doneList       list.Model
	log            string
	focusedSection Status
	screenWidth    int
	screenHeight   int
}

func newModel() Model {
	m := Model{}

	m.todoItems = map[Status][]TodoItem{
		todo:       make([]TodoItem, 0),
		inProgress: make([]TodoItem, 0),
		done:       make([]TodoItem, 0),
	}

	m.todoItems[todo] = []TodoItem{
		{title: "one", description: "One", status: todo},
		{title: "two", description: "two", status: todo},
	}
	m.todoItems[inProgress] = []TodoItem{
		{title: "three", description: "three", status: inProgress},
	}
	m.todoItems[done] = []TodoItem{
		{title: "four", description: "four", status: done},
	}

	m.todoList = list.New()
	m.todoList.Title = "Todo"
    m.todoList.Focus()
	m.inProgressList = list.New()
	m.inProgressList.Title = "In Progress"
	m.doneList = list.New()
	m.doneList.Title = "Done"

	for _, item := range m.todoItems[todo] {
		m.todoList.AddItem(item)
	}
	for _, item := range m.todoItems[inProgress] {
		m.inProgressList.AddItem(item)
	}
	for _, item := range m.todoItems[done] {
		m.doneList.AddItem(item)
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) onAddedItem(title string) {
	newItem := TodoItem{title: title, description: "", status: m.focusedSection}

	var focusedList *list.Model
	switch m.focusedSection {
	case todo:
		focusedList = &m.todoList
	case inProgress:
		focusedList = &m.inProgressList
	case done:
		focusedList = &m.doneList
	}

	focusedItems := m.todoItems[m.focusedSection]

	idx := focusedList.SelectedIndex()
	head := focusedItems[:idx]
	tailSource := focusedItems[idx:]
	tail := make([]TodoItem, len(tailSource))
	copy(tail, tailSource)
	head = append(head, newItem)
	m.todoItems[m.focusedSection] = append(head, tail...)

	listItems := make([]list.ListItem, len(m.todoItems[m.focusedSection]))
	for idx := range m.todoItems[m.focusedSection] {
		listItems[idx] = m.todoItems[m.focusedSection][idx]
	}
	focusedList.SetItems(listItems...)
}

func (m *Model) onFocusNextList() {
    switch m.focusedSection {
    case todo:
        m.focusedSection = inProgress
        m.todoList.Blur()
        m.inProgressList.Focus()
    case inProgress:
        m.focusedSection = done
        m.inProgressList.Blur()
        m.doneList.Focus()
    case done:
        m.focusedSection = todo
        m.doneList.Blur()
        m.todoList.Focus()
    }
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.screenWidth = msg.Width
		m.screenHeight = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
        case "tab":
            m.onFocusNextList()
		}
	case list.AddedItem:
		m.onAddedItem(msg.Value)
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.focusedSection == todo {
		m.todoList, cmd = m.todoList.Update(msg)
		cmds = append(cmds, cmd)
	}
	if m.focusedSection == inProgress {
		m.inProgressList, cmd = m.inProgressList.Update(msg)
		cmds = append(cmds, cmd)
	}
	if m.focusedSection == done {
		m.doneList, cmd = m.doneList.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)

}

func (m Model) View() string {
	listwidth := (m.screenWidth - 12) / 3
	listheight := m.screenHeight - 6
	focusedListStyle = focusedListStyle.Width(listwidth).Height(listheight)
	baseListStyle = baseListStyle.Width(listwidth).Height(listheight)

	todoList := m.todoList.View()
	inProgressList := m.inProgressList.View()
	doneList := m.doneList.View()

	switch m.focusedSection {
	case todo:
		todoList = focusedListStyle.Render(todoList)
		inProgressList = baseListStyle.Render(inProgressList)
		doneList = baseListStyle.Render(doneList)
	case inProgress:
		todoList = baseListStyle.Render(todoList)
		inProgressList = focusedListStyle.Render(inProgressList)
		doneList = baseListStyle.Render(doneList)
	case done:
		todoList = baseListStyle.Render(todoList)
		inProgressList = baseListStyle.Render(inProgressList)
		doneList = focusedListStyle.Render(doneList)
	}

	list := lipgloss.JoinHorizontal(
		lipgloss.Center,
		todoList,
		inProgressList,
		doneList,
	)
	if m.log != "" {
		list += "\n\n" + "Log: " + m.log
	}
	return list + "\n\n"
}

func main() {
	model := newModel()
	p := tea.NewProgram(model, tea.WithAltScreen())

	_, err := p.Run()
	if err != nil {
		panic(err)
	}
}
