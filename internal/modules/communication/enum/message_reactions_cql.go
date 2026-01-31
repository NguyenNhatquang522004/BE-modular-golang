package enum

import (
	"fmt"

	"github.com/gocql/gocql"
)

// MarshalCQL: Go (Int) -> Cassandra (String)
// VD: ReactionHeart (0) -> "heart"
func (e MessageReactionCode) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return []byte(e.String()), nil
}

// UnmarshalCQL: Cassandra (String) -> Go (Int)
// VD: "heart" -> ReactionHeart (0)
func (e *MessageReactionCode) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	val, err := MessageReactionCodeString(string(data))
	if err != nil {
		// Fallback: Nếu DB có mã lạ (emoji lạ), có thể log warning hoặc trả về default
		return fmt.Errorf("unknown reaction code: %s", string(data))
	}
	*e = val
	return nil
}