package tasmota

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

const (
	MAX_COMMAND_LENGTH = 30
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
	ErrNoSubscribers         = fmt.Errorf("no subscribers")
	ErrTooMany               = fmt.Errorf("too many commands")
)

type Client struct {
	Type       ClientType
	Addr       string
	NoDelay    bool
	MQTTConfig MQTTConfig
	cm         *autopaho.ConnectionManager
}

type MQTTConfig struct {
	BrokerUrl         string
	Topic             string
	Username          string
	Password          []byte
	ClientID          string
	KeepAlive         uint16
	ConnectRetryDelay time.Duration
	Timeout           time.Duration
}

func (t *Client) init(ctx context.Context) error {
	var (
		u   *url.URL
		err error
	)
	if t.Type == ClientTypeMQTT && t.cm == nil {
		u, err = url.Parse(t.MQTTConfig.BrokerUrl)
		if err != nil {
			return err
		}
		cliCfg := autopaho.ClientConfig{
			BrokerUrls:        []*url.URL{u},
			KeepAlive:         t.MQTTConfig.KeepAlive,
			ConnectRetryDelay: t.MQTTConfig.ConnectRetryDelay,
			ConnectTimeout:    t.MQTTConfig.Timeout,
			ClientConfig: paho.ClientConfig{
				ClientID: t.MQTTConfig.ClientID,
			},
		}
		cliCfg.SetUsernamePassword(t.MQTTConfig.Username, t.MQTTConfig.Password)
		t.cm, err = autopaho.NewConnection(ctx, cliCfg)
	}
	return err
}

func (t Client) Command() Command {
	return Command{
		Client: t,
	}
}

func (t Client) Execute(ctx context.Context, c Command) error {
	var (
		cmd string
		err error
		pr  *paho.PublishResponse
	)

	err = t.init(ctx)
	if err != nil {
		return err
	}

	switch t.Type {
	case ClientTypeWeb:
		cmd, err = t.Build(c)
		if err != nil {
			return err
		}
		_, err = http.Get(fmt.Sprintf("http://%s/cm?cmnd=%s", t.Addr, url.QueryEscape(cmd)))
		return err
	case ClientTypeMQTT:
		var payload string

		err = t.cm.AwaitConnection(ctx)
		if err != nil {
			return err
		}
		if len(c.commands) > 1 {
			payload, err = t.Build(c)
			if err != nil {
				return err
			}
			if t.NoDelay {
				cmd = "backlog0"
			} else {
				cmd = "backlog"
			}
		} else if len(c.commands) == 1 {
			cmd = c.commands[0][0]
			payload = c.commands[0][1]
		}

		pr, err = t.cm.Publish(ctx, &paho.Publish{
			QoS:     2,
			Topic:   fmt.Sprintf("cmnd/%s/%s", t.MQTTConfig.Topic, cmd),
			Payload: []byte(payload),
		})
		if err != nil {
			return err
		} else if pr.ReasonCode != 0 && pr.ReasonCode != 16 {
			// 16 = Server received message but there are no subscribers
			return ErrNoSubscribers
		}

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
		cmd    string
	)

	if len(c.commands) > MAX_COMMAND_LENGTH {
		return "", ErrTooMany
	}

	for _, cmd := range c.commands {
		res = append(res, fmt.Sprintf("%s %s", cmd[0], cmd[1]))
	}

	cmd = strings.Join(res, ";")

	switch t.Type {
	case ClientTypeWeb:
		if len(c.commands) > 0 {
			if c.NoDelay {
				prefix = "backlog0 "
			} else {
				prefix = "backlog "
			}
		}

		return prefix + cmd, nil
	case ClientTypeTest:
		fallthrough
	case ClientTypeMQTT:
		return cmd, nil
	default:
		return "", ErrUnsupportedClientType
	}
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
