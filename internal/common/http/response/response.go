package response

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
type OptionalResponse func(*Response)

func WithData(data interface{}) OptionalResponse {
	return func(r *Response) {
		r.Data = data
	}
}
func WithMessage(message string) OptionalResponse {
	return func(r *Response) {
		r.Message = message
	}
}
func WithStatus(status string) OptionalResponse {
	return func(r *Response) {
		r.Status = status
	}
}

func SuccessResponse(message string, data interface{}) Response {
	return Response{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(message string) Response {
	return Response{
		Status:  "error",
		Message: message,
	}
}
func NewResponse(opts ...OptionalResponse) Response {
	res := Response{
		Status:  "success",
		Message: "",
		Data:    nil,
	}
	for _, opt := range opts {
		opt(&res)
	}
	return res
}
