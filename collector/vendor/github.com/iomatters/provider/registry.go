package provider

import (
	"fmt"
)

type cryptocompareResponse struct {
	Raw     map[string]map[string]map[string]interface{} `json:"RAW"`
	Display map[string]map[string]map[string]string      `json:"DISPLAY"`
}

type Provider interface {
	GetAuthorized() error
	Pull(fsyms []string, tsyms []string) (*cryptocompareResponse, error)
}

func NewProvider(id string) (Provider, error) {
	switch id {
	case "cryptocompare":
		return NewCryptocompareProvider()
	default:
		return nil, fmt.Errorf("invalid provider ID: %s", id)
	}

}
