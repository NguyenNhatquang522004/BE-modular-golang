package sharedEnums

type ProcessingStatus int

const (
	ProcessingPending    ProcessingStatus = iota // 'pending' (Mới upload, chưa xử lý)
	ProcessingProcessing                         // 'processing' (Đang transcode/nén)
	ProcessingActive                             // 'active' (Đã xong, user có thể xem)
	ProcessingFailed                             // 'failed' (Lỗi file)
	ProcessingPaused                             // 'deleted' (Người dùng xóa video)
)
