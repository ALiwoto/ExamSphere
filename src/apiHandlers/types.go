package apiHandlers

type EndpointResponse struct {
	Success bool           `json:"success"`
	Result  any            `json:"result"`
	Error   *EndpointError `json:"error"`
}

type EndpointError struct {
	ErrorCode int    `json:"code"`
	Message   string `json:"message"`
	Origin    string `json:"origin"`
	Date      string `json:"date"`
}
