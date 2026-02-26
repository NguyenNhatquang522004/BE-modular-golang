package enum

// =============================================================================
// 1. COMMENT STATUS
// =============================================================================

//go:generate enumer -type=CommentStatus -json -transform=snake -trimprefix=CommentStatus
type CommentStatus int

const (
	CommentStatusActive  CommentStatus = iota // 'active'
	CommentStatusHidden                       // 'hidden'
	CommentStatusDeleted                      // 'deleted'
	CommentStatusPending                      // 'pending' (Chờ duyệt nếu có bật mode kiểm duyệt)
)

// =============================================================================
// 2. COMMENT MEDIA TYPE
// =============================================================================

//go:generate enumer -type=CommentMediaType -json -transform=snake -trimprefix=CommentMedia