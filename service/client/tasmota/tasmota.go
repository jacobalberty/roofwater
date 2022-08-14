package tasmota

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ClientType int

const (
	ClientTypeWeb ClientType = iota
	ClientTypeMQTT
	ClientTypeSerial
	ClientTypeTest
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
	Type    ClientType
	Addr    string
	NoDelay bool
}

func (t Client) Command() Command {
	return Command{
		Client: t,
	}
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
		_, err := http.Get(fmt.Sprintf("http://%s/cm?cmnd=%s", t.Addr, url.QueryEscape(cmd)))
		return err
	case ClientTypeTest:
		log.Println(cmd)
		return nil
	default:
		return ErrUnsupportedClientType
	}
}

func (t Client) Build(c Command) (string, error) {
	var (
		res    = make([]string, 0, len(c.commands))
		prefix string
	)
	if len(c.commands) > 0 {
		if c.NoDelay {
			prefix = "backlog0 "
		} else {
			prefix = "backlog "
		}
	}
	for _, cmd := range c.commands {
		res = append(res, fmt.Sprintf("%s %s", cmd[0], cmd[1]))
	}

	return prefix + strings.Join(res, ";"), nil
}

type Command struct {
	Client
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
