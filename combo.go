package main

import (
	"errors"
	"strings"
)

type Combo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Raw      string `json:"raw"`
}

func (c *Combo) ToString() string {
	return c.Username + ":" + c.Password
}

func NewComboFromString(raw string) (*Combo, error) {

	if len(raw) == 0 || !strings.Contains(raw, ":") {
		return nil, errors.New("invalid combo")
	}

	split := strings.Split(raw, ":")

	return &Combo{
		Username: split[0],
		Password: split[1],
		Raw:      raw,
	}, nil
}

func NewCombo(username, password, raw string) (*Combo, error) {

	if len(username) == 0 {
		return nil, errors.New("invalid combo")
	}

	return &Combo{
		Username: username,
		Password: password,
		Raw:      raw,
	}, nil
}
