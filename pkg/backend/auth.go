package backend

import (
	"errors"
	"net"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/mdns"
)

// TODO: check if auth in config is valid
func IsLoggedIn() bool {
	return false
}

type IpMsg struct{ net.IP }
type ErrMsg struct{ Err error }

func (e ErrMsg) Error() string { return e.Err.Error() }

func GetIp() tea.Msg {
	entriesCh := make(chan *mdns.ServiceEntry)
	ipCh := make(chan net.IP)
	defer close(entriesCh)
	defer close(ipCh)

	go func() {
		for entry := range entriesCh {
			ipCh <- entry.AddrV4
		}
	}()

	if err := mdns.Lookup("_hue._tcp", entriesCh); err != nil {
		return ErrMsg{errors.New("LOOKUP ERROR")}
	}

	select {
	case ip := <-ipCh:
		return IpMsg{ip}
	case <-time.After(1 * time.Second):
		return ErrMsg{errors.New("No IP adress found")}
	}
}
