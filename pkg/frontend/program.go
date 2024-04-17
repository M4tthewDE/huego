package frontend

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/m4tthewde/huego/pkg/backend"
	"net"
)

func NewProgram() *tea.Program {
	return tea.NewProgram(initialModel())
}

type Model struct {
	Err        error
	loggedIn   bool
	fetchingIp bool
	ip         net.IP
	spinner    spinner.Model
}

func initialModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return Model{
		loggedIn:   backend.IsLoggedIn(),
		fetchingIp: false,
		spinner:    s,
	}
}

func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case backend.IpMsg:
		m.ip = msg.IP
		return m, tea.Quit
	case backend.ErrMsg:
		m.Err = msg
		return m, tea.Quit
	default:
		if !m.loggedIn {
			if !m.fetchingIp {
				m.fetchingIp = true
				return m, backend.GetIp
			}
		}

		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	if !m.loggedIn {
		if m.ip == nil {
			return m.spinner.View() + "Fetching IP adress of your Hue Bridge"
		} else {
			return fmt.Sprintf("%s Authenticating with Hue Bridge at %v", m.spinner.View(), m.ip)
		}
	} else {
		return "TODO: lamps"
	}
}
