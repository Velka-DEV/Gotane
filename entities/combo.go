package entities

type Combo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Raw      string `json:"raw"`
}

func NewCombo(username, password, raw string) *Combo {
	return &Combo{
		Username: username,
		Password: password,
		Raw:      raw,
	}
}
