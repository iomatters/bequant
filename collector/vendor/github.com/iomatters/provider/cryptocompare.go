package provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

const (
	baseURL     = "https://min-api.cryptocompare.com"
	api         = "data/pricemultifull"
	contentType = "application/json"
	httpTimeout = time.Second * 5
)

type CryptocompareProvider struct {
	BaseURL     string
	API         string
	ContentType string
	HttpTimeout time.Duration
}

func NewCryptocompareProvider() (*CryptocompareProvider, error) {
	return &CryptocompareProvider{
		BaseURL:     baseURL,
		API:         api,
		ContentType: contentType,
		HttpTimeout: httpTimeout,
	}, nil
}

func (p *CryptocompareProvider) GetAuthorized() error {
	return nil
}

func (p *CryptocompareProvider) Pull(fsyms []string, tsyms []string) (*cryptocompareResponse, error) {
	url := fmt.Sprintf("%s/%s?fsyms=%s&tsyms=%s", p.BaseURL, p.API, strings.Join(fsyms, ","), strings.Join(tsyms, ","))
	client := http.Client{
		Timeout: p.HttpTimeout,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", p.ContentType)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var unpacked cryptocompareResponse
	if err := json.Unmarshal(body, &unpacked); err != nil {
		return nil, err
	}

	// Suspecting error, return RAW response
	if reflect.DeepEqual(unpacked, cryptocompareResponse{}) {
		return nil, fmt.Errorf("%s", string(body))
	}

	return &unpacked, nil
}
