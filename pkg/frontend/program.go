package frontend

import (
	"fmt"
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
}

func initialModel() Model {
	return Model{
		loggedIn:   backend.IsLoggedIn(),
		fetchingIp: false,
	}
}

func (m Model) Init() tea.Cmd {
	if !m.loggedIn {
		return backend.GetIp
	} else {
		return nil
	}
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
		m.fetchingIp = false
		return m, nil
	case backend.ErrMsg:
		m.Err = msg.Err
		return m, tea.Quit

	default:
		if !m.loggedIn {
			if !m.fetchingIp {
				m.fetchingIp = true
				return m, backend.GetIp
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	if !m.loggedIn {
		if m.ip == nil {
			return "Fetching IP adress of Hue Bridge"
		} else {
			return fmt.Sprintf("Authenticating with Hue Bridge %v", m.ip)
		}
	} else {
		return "TODO: lamps"
	}
}
