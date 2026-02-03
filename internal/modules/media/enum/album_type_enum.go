package enum

//go:generate enumer -type=AlbumType -json -transform=snake -trimprefix=AlbumType
type AlbumType int

const (
	AlbumTypeNormal  AlbumType = iota // 'normal' (Album thường do user tạo)
	AlbumTypeProfile                  // 'profile' (Chứa ảnh đại diện)
	AlbumTypeCover                    // 'cover' (Chứa ảnh bìa)
	AlbumTypeMobile                   // 'mobile' (Ảnh tải lên từ điện thoại)
	AlbumTypeTrash                    // 'trash' (Thùng rác - Soft delete logic)
)