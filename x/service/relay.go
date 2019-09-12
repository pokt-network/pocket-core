package service

import (
	"bytes"
	"errors"
	"github.com/pokt-network/pocket-core/types"
	"io/ioutil"
	"net/http"
)

type Relay struct {
	Blockchain       ServiceBlockchain `json:"blockchain"`
	Payload          ServicePayload    `json:"payload"`
	ServiceToken     ServiceToken      `json:"token"`
	IncrementCounter IncrementCounter  `json:"incrementCounter"`
}

func (r Relay) Validate(hostedBlockchains ServiceBlockchains) error {
	if r.Blockchain == nil || len(r.Blockchain) == 0 {
		return EmptyBlockchainError
	}
	if r.Payload.Data.Bytes() == nil || len(r.Payload.Data.Bytes()) == 0 {
		return EmptyPayloadDataError
	}
	if err := r.ServiceToken.Validate(); err != nil {
		return errors.New(InvalidTokenError.Error() + " : " + err.Error())
	}
	if err := r.IncrementCounter.Validate(r.ServiceToken.AATMessage.ClientPublicKey, r.ServiceToken.Hash()); err != nil {
		return errors.New(InvalidIncrementCounterError.Error() + " : " + err.Error())
	}
	if !types.Blockchains(hostedBlockchains).Contains(r.Blockchain.String()) {
		return UnsupportedBlockchainError
	}
	if r.Payload.HttpServicePayload != (HttpServicePayload{}) {
		if len((r.Payload).Method) == 0 {
			r.Payload.Method = DEFAULTHTTPMETHOD
		}
	}
	return nil
}

func (r Relay) Execute(hostedBlockchains ServiceBlockchains) (string, error) {
	if err := r.Validate(hostedBlockchains); err != nil {
		return "", err
	}
	chainURL, err := hostedBlockchains.GetChainURL(r.Blockchain.String())
	if err != nil {
		return "", err
	}
	result, err := executeHTTPRequest(r.Payload.Data.Bytes(), chainURL, r.Payload.Method)
	return result, err
}

// "executeHTTPRequest" takes in the raw json string and forwards it to the HTTP endpoint
func executeHTTPRequest(payload []byte, url string, method string) (string, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", err
	}
	// todo inspect unhandled error
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}
