package tasmota

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ClientType int

const (
	ClientTypeWeb ClientType = iota
	ClientTypeMQTT
)

type PowerState int

const (
	PowerOff PowerState = iota
	PowerOn
	PowerToggle
)

var (
	ErrUnsupportedClientType = fmt.Errorf("unsupported client type")
)

type Client struct {
	Type ClientType
	IP   net.IP
}

func (t Client) Execute(c Command) error {
	var (
		cmd string
		err error
	)
	cmd, err = t.Build(c)
	if err != nil {
		return err
	}
	switch t.Type {
	case ClientTypeWeb:
		_, err := http.Get(fmt.Sprintf("http://%s/cm?cmnd=%s", t.IP.String(), cmd))
		return err
	}
	return nil
}

func (t Client) Build(c Command) (string, error) {
	var (
		res    = make([]string, len(c.commands))
		prefix string
	)
	if len(c.commands) > 0 {
		prefix = "backlog "
	}
	for _, cmd := range c.commands {
		res = append(res, fmt.Sprintf("%s %s\n", cmd[0], cmd[1]))
	}
	switch t.Type {
	case ClientTypeWeb:
		return url.QueryEscape(prefix + strings.Join(res, ";")), nil
	case ClientTypeMQTT:
		return prefix + strings.Join(res, ";"), nil
	}
	return "", ErrUnsupportedClientType
}

type Command struct {
	commands [][]string
}

func (t Command) Power(state PowerState) Command {
	t.commands = append(t.commands, []string{"power", fmt.Sprint(state)})
	return t
}

func (t Command) Delay(d time.Duration) Command {
	t.commands = append(t.commands, []string{"delay", fmt.Sprint(d.Milliseconds() / 100)})
	return t
}
