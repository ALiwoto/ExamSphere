package apiHandlers

// APIErrorCode is an enum that represents the error codes that can be
// returned by the API.
type APIErrorCode int

// EndpointResponse is the structure of the response that is sent to the
// client.
type EndpointResponse struct {
	Success bool           `json:"success"`
	Result  any            `json:"result"`
	Error   *EndpointError `json:"error"`
}

// EndpointError is the structure of the error that is sent to the client.
type EndpointError struct {
	ErrorCode APIErrorCode `json:"code"`
	Message   string       `json:"message"`
	Origin    string       `json:"origin"`
	Date      string       `json:"date"`
}
