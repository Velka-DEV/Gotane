package pkg

import (
	"errors"
	"strings"
)

type Combo struct {
	Username string
	Password string
	Raw      string
}

func ParseCombo(s string) (*Combo, error) {
	parts := strings.Split(s, ":")

	if len(parts) != 2 {
		return &Combo{}, errors.New("invalid combo format")
	}

	return &Combo{
		Username: strings.TrimSpace(parts[0]),
		Password: strings.TrimSpace(parts[1]),
		Raw:      s,
	}, nil
}

func (c Combo) String() string {
	return c.Username + ":" + c.Password
}

func (c Combo) IsValid() bool {
	return c.Username != "" && c.Password != ""
}
