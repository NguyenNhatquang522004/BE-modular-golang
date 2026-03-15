package dto

// OMEWebhookPayload đại diện cho payload JSON từ OvenMediaEngine gửi sang
type OMEWebhookPayload struct {
	Request struct {
		Direction string `json:"direction"`
		Protocol  string `json:"protocol"`
	} `json:"request"`
	App struct {
		Name string `json:"name"`
	} `json:"app"`
	Stream struct {
		Name string `json:"name"` // Đây chính là LiveSessionID của bạn
	} `json:"stream"`
	Event struct {
		Name string `json:"name"` // "stream.created" hoặc "stream.destroyed"
	} `json:"event"`
}
