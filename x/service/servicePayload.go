package service

type ServicePayload struct {
	Data ServiceData `json:"data"` // the payload data for the non native chain
	HttpServicePayload
}

type HttpServicePayload struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

type ServiceData string

func (sd ServiceData) Bytes() []byte {
	// todo
	return []byte(sd)
}
