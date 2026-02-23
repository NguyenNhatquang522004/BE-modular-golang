package mongodbErrors

type BulkError struct {
	ID     string
	Reason string
}

type EditLogsBulkError struct {
	ID       string `json:"id"`        // ID của document bị lỗi
	TargetID string `json:"target_id"` // TargetID của document bị lỗi
	EditorID string `json:"editor_id"` // EditorID của document bị lỗi
	Reason   string `json:"reason"`    // Chi tiết lỗi
}
