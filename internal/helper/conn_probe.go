package helper

import (
	logger "golang-api-starter/internal/helper/logger/zap_log"
	"net"
	"time"
)

type ProbeTarget struct {
	Host        string        // IP or hostname
	Ports       []string      // list of ports to be tested
	NetProtocol string        // TCP / UDP
	TimeoutSec  time.Duration // time.Second for test the port
	DialFunc    IDial
}

type IDial interface {
	Probe(string) error
}

func (pt *ProbeTarget) Probe(port string) error {
	_, err := net.DialTimeout(pt.NetProtocol, net.JoinHostPort(pt.Host, port), pt.TimeoutSec)
	return err
}

func (pt *ProbeTarget) PortsProbe() error {
	for _, port := range pt.Ports {
		err := pt.DialFunc.Probe(port)
		if err != nil {
			return logger.Errorf("Connection error: %s, err: %+s", net.JoinHostPort(pt.Host, port), err.Error())
		}
	}
	return nil
}
