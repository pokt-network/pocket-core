package service

import "reflect"

type ServicePayload struct {
	Data ServiceData `json:"data"` // the payload data for the non native chain
	HttpServicePayload
}

type HttpServicePayload struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

type ServiceData string

type ServiceType int

func (sd ServiceData) Bytes() []byte {
	// todo
	return []byte(sd)
}

func (sp ServicePayload) Type() ServiceType {
	if reflect.TypeOf(sp.HttpServicePayload) == reflect.TypeOf((HttpServicePayload{})) {
		return REST
	}
	return RPC
}
