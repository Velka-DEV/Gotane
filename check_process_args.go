package gotane

import (
	"net/http"
)

type CheckProcessArgs struct {
	Combo *Combo `json:"combo"`

	// Proxy should not be used to create new client (this will decrease performances a lot)
	// Except in case if a custom client is required
	Proxy *Proxy `json:"proxy"`

	Client *http.Client `json:"client"`
}
