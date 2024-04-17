package backend

import (
	"errors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/mdns"
	"net"
	"time"
)

// TODO: check if auth in config is valid
func IsLoggedIn() bool {
	return false
}

type IpMsg struct{ net.IP }
type ErrMsg struct{ err error }

func (e ErrMsg) Error() string { return e.err.Error() }

func GetIp() tea.Msg {
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	ipCh := make(chan net.IP, 4)
	defer close(entriesCh)
	defer close(ipCh)

	go func() {
		for entry := range entriesCh {
			ipCh <- entry.AddrV4
		}
	}()

	if err := mdns.Lookup("_hue._tcp.local.", entriesCh); err != nil {
		return ErrMsg{err}
	}

	select {
	case ip := <-ipCh:
		return IpMsg{ip}
	case <-time.After(1 * time.Second):
		return ErrMsg{errors.New("No IP adress found")}
	}
}
