package helper

import (
	"errors"
	"fmt"
	zlog "golang-api-starter/internal/helper/logger/zap_log"
	"net"
	"strconv"
	"testing"
	"time"
)

type probeTest struct {
	name  string
	input *ProbeTarget
	want  error
}

type PortSuccess struct{}

func (ps *PortSuccess) Probe(string) error {
	return nil
}

type PortFailed struct{}

func (ps *PortFailed) Probe(string) error {
	return errors.New("Failed connect to port")
}

func TestPortProbe(t *testing.T) {
	cfg.LoadEnvVariables()
	zlog.NewZlog()

	probeTests := []probeTest{
		{
			name: "test port ok",
			input: &ProbeTarget{
				Host:        cfg.Notification.Smtp.Host,
				Ports:       []string{strconv.Itoa(cfg.Notification.Smtp.Port)},
				NetProtocol: "tcp",
				TimeoutSec:  3 * time.Second,
				DialFunc:    &PortSuccess{},
			},
			want: nil,
		},
		{
			name: "test port failed",
			input: &ProbeTarget{
				Host:        cfg.Notification.Smtp.Host,
				Ports:       []string{strconv.Itoa(cfg.Notification.Smtp.Port)},
				NetProtocol: "tcp",
				TimeoutSec:  3 * time.Second,
				DialFunc:    &PortFailed{},
			},
			want: errors.New(fmt.Sprintf("Connection error: %s", net.JoinHostPort(cfg.Notification.Smtp.Host, strconv.Itoa(cfg.Notification.Smtp.Port)))),
		},
	}

	for _, testCase := range probeTests {
		got := testCase.input.PortsProbe()
		var eq bool = true
		if got != nil {
			eq = testCase.want.Error() == got.Error()
		} 

		if !eq {
			t.Errorf("got %q want %q", got, testCase.want)
		}
	}
}
