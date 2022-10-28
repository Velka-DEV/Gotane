package main

import (
	"errors"
	"strconv"
	"strings"
)

type Proxy struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Retries  int    `json:"retries"`
	IsBanned bool   `json:"banned"`
}

func NewProxy(protocol string, host string, port int) (*Proxy, error) {

	if len(protocol) == 0 || len(host) == 0 || port == 0 {
		return nil, errors.New("invalid proxy")
	}

	return &Proxy{
		Protocol: protocol,
		Host:     host,
		Port:     port,
		Retries:  0,
		IsBanned: false,
	}, nil
}

func NewProxyFromUrlString(raw string) (*Proxy, error) {

	if len(raw) == 0 || !strings.Contains(raw, "://") || !strings.Contains(raw, ":") {
		return nil, errors.New("invalid proxy")
	}

	split := strings.Split(raw, "://")
	protocol := split[0]

	split = strings.Split(split[1], ":")
	host := split[0]
	port, err := strconv.Atoi(split[1])

	if err != nil {
		return nil, err
	}

	return NewProxy(protocol, host, port)
}

func NewProxyFromString(raw string, protocol string) (*Proxy, error) {

	if len(raw) == 0 || !strings.Contains(raw, ":") {
		return nil, errors.New("invalid proxy")
	}

	split := strings.Split(raw, ":")
	host := split[0]
	port, err := strconv.Atoi(split[1])

	if err != nil {
		return nil, err
	}

	return NewProxy(protocol, host, port)
}

func (p *Proxy) ToString() string {
	return p.Protocol + "://" + p.Host + ":" + strconv.Itoa(p.Port)
}
