package enum

//go:generate enumer -type=PrivacySearch -json  -yaml -trimprefix=Privacy_Search_ -transform=snake

type PrivacySearch int

const (
	// PrivacyUnknown là giá trị mặc định (0) để tránh lỗi zero-value
	Privacy_Search_Unknown PrivacySearch = iota

	// PrivacyPublic: Nhóm công khai
	Privacy_Search_Public

	// PrivacyClosed: Nhóm kín (nhưng vẫn có thể index tên)
	Privacy_Search_Closed

	// PrivacySecret: Nhóm bí mật (Không index)
	Privacy_Search_Secret
)
