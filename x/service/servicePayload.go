package service

type ServicePayload struct {
	Data   ServiceData `json:"data"`   // the payload data for the non native chain
	Method string      `json:"method"` // the http method needed for the rest call
	Path   string      `json:"path"`   // the path needed for REST calls
}

type ServiceData string

func (sd ServiceData) Bytes() []byte {
	// todo
	return []byte(sd)
}
