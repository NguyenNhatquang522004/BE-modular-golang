package enum

// =============================================================================
// PAGE ROLE (RBAC)
// =============================================================================

//go:generate enumer -type=PageRole -json -transform=snake -trimprefix=Role
type PageRole int

const (
	RoleAdmin      PageRole = iota // 'admin': Toàn quyền (Quản lý role, xóa page...)
	RoleEditor                     // 'editor': Đăng bài, sửa bài, xem insight
	RoleModerator                  // 'moderator': Trả lời inbox, comment, ban user
	RoleAdvertiser                 // 'advertiser': Chỉ tạo quảng cáo, xem insight
	RoleAnalyst                    // 'analyst': Chỉ xem chỉ số (Insight)
)