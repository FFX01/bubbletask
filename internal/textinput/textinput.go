package textinput

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const cursor string = "█"

var (
	cursorBackround  lipgloss.Color = lipgloss.Color("1")
	cursorForeground lipgloss.Color = lipgloss.Color("2")
	cursorStyle      lipgloss.Style = lipgloss.NewStyle().
				Background(cursorBackround).
				Foreground(cursorForeground)
)

type Model struct {
	buf       []rune
	focused   bool
	cursoridx int
}

func New() Model {
    m := Model{}
    m.buf = make([]rune, 0)
    return m
}

func (m *Model) Focus() {
	m.focused = true
}

func (m *Model) Blur() {
	m.focused = false
}

func (m *Model) Reset() {
    m.buf = []rune{}
    m.focused = false
    m.cursoridx = 0
}

func (m *Model) SetValue(v string) {
	m.buf = []rune(v)
	if m.cursoridx >= len(m.buf) {
		m.cursoridx = len(m.buf)
	}
}

func (m *Model) Value() string {
	return string(m.buf)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) insertRunes(v []rune) {
	head := m.buf[:m.cursoridx]
	tailSource := m.buf[m.cursoridx:]
	tail := make([]rune, len(tailSource))
	copy(tail, tailSource)

	for _, r := range v {
		head = append(head, r)
		m.cursoridx++
	}

	m.buf = append(head, tail...)
}

func (m *Model) onLeft() {
	if m.cursoridx > 0 {
		m.cursoridx--
	}
}

func (m *Model) onRight() {
	if m.cursoridx < len(m.buf) {
		m.cursoridx++
	}
}

func (m *Model) onBackspace() {
	if m.cursoridx < 1 {
		return
	}

	if m.cursoridx == len(m.buf) {
		m.buf = m.buf[:m.cursoridx-1]
		m.cursoridx--
		return
	}

	m.buf = append(m.buf[:m.cursoridx-1], m.buf[m.cursoridx:]...)
	m.cursoridx--
}

func (m *Model) onDelete() {
	if m.cursoridx == len(m.buf) || len(m.buf) < 1 {
		return
	}

	m.buf = append(m.buf[:m.cursoridx], m.buf[m.cursoridx+1:]...)
}

type Confirm struct {
	Value string
}

func (m *Model) confirm() tea.Msg {
	return Confirm{
		Value: m.Value(),
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    if !m.focused {
        return m, nil
    }

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, m.confirm
		case "left":
			m.onLeft()
		case "right":
			m.onRight()
		case "backspace":
			m.onBackspace()
		case "delete":
			m.onDelete()
		case " ":
			m.insertRunes([]rune{' '})
        case "home":
            m.cursoridx = 0
        case "end":
            m.cursoridx = len(m.buf) - 1
		default:
			if msg.Type == tea.KeyRunes {
				m.insertRunes(msg.Runes)
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	if len(m.buf) < 1 {
		return cursorStyle.Render(" ")
	}
	if m.cursoridx >= len(m.buf) {
		return string(m.buf) + cursorStyle.Render(" ")
	}

	atCursor := string(m.buf[m.cursoridx])
	beforeCursor := string(m.buf[:m.cursoridx])
	afterCursor := string(m.buf[m.cursoridx+1:])

	styledCursorPos := cursorStyle.Render(atCursor)

	output := beforeCursor + styledCursorPos + afterCursor

	return output
}
