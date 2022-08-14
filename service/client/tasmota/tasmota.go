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
)

type Client struct {
	Type       ClientType
	Addr       string
	NoDelay    bool
	MQTTConfig MQTTConfig
	cm         *autopaho.ConnectionManager
}

type MQTTConfig struct {
	BrokerUrl string
	Topic     string
	Username  string
	Password  []byte
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
			BrokerUrls: []*url.URL{u},
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

	cmd, err = t.Build(c)
	if err != nil {
		return err
	}

	switch t.Type {
	case ClientTypeWeb:
		_, err = http.Get(fmt.Sprintf("http://%s/cm?cmnd=%s", t.Addr, url.QueryEscape(cmd)))
		return err
	case ClientTypeMQTT:
		err = t.cm.AwaitConnection(ctx)
		if err != nil {
			return err
		}

		pr, err = t.cm.Publish(ctx, &paho.Publish{
			Topic:   t.MQTTConfig.Topic,
			Payload: []byte(cmd),
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
