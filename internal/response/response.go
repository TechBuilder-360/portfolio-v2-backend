package response

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func ErrorResponse(message string) Response {
	return Response{
		Status:  false,
		Message: message,
	}
}

func SuccessResponse(message string, data, meta interface{}) Response {
	return Response{
		Status:  true,
		Message: message,
		Data:    data,
		Meta:    meta,
	}
}
