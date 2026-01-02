package apis

type errorDataResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func ErrorDataResponse(err string, message string, status int) *errorDataResponse {
	return &errorDataResponse{
		Error:   err,
		Message: message,
		Status:  status,
	}
}
