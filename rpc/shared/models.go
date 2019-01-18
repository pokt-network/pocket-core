package shared

// "JSONResponse" is a metadata and data response in JSON format.
type JSONResponse struct {
	Data string `json:"data"`
}

// "JSONErrorResponse" is an error response in JSON format.
type JSONErrorResponse struct {
	Error *APIError `json:"error"`
}

// "APIError" is an error feedback structure containing a title and a status.
type APIError struct {
	Status int    `json:"error"`
	Title  string `json:"title"`
}

// "APIReference' is an in-client API reference.
type APIReference struct {
	Endpoint   string   `json:"endpoint"`
	Method     string   `json:"method"`
	Parameters []string `json:"params"`
	Returns    string   `json:"returns"`
	Example    string   `json:"example"`
}
