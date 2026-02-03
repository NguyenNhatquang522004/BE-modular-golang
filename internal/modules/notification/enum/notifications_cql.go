package enum

import (
	"fmt"

	"github.com/gocql/gocql"
)

// MarshalCQL: Go (Int) -> Cassandra (String)
// VD: NotifPostLike -> "post_like"
func (e NotificationType) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return []byte(e.String()), nil
}

// UnmarshalCQL: Cassandra (String) -> Go (Int)
// VD: "post_like" -> NotifPostLike
func (e *NotificationType) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	val, err := NotificationTypeString(string(data))
	if err != nil {
		// Fallback: Nếu có loại thông báo mới chưa update code, trả về lỗi hoặc log warning
		return fmt.Errorf("unknown notification type: %s", string(data))
	}
	*e = val
	return nil
}