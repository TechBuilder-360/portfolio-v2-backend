package response

type SuccessResp struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type SuccessMetaResp struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type ValidationErrorRes struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

func ErrorResponse(message string) Response {
	return Response{
		Status:  false,
		Message: message,
	}
}

func ValidationErrorResponse(data interface{}) ValidationErrorRes {
	return ValidationErrorRes{
		Status:  false,
		Message: "validation failed",
		Errors:  data,
	}
}

func SuccessResponse(message string) Response {
	return Response{
		Status:  true,
		Message: message,
	}
}

func DataResponse(message string, data interface{}) SuccessResp {
	return SuccessResp{
		Status:  true,
		Message: message,
		Data:    data,
	}
}

func MetaResponse(message string, data, meta interface{}) SuccessMetaResp {
	return SuccessMetaResp{
		Status:  true,
		Message: message,
		Data:    data,
		Meta:    meta,
	}
}
