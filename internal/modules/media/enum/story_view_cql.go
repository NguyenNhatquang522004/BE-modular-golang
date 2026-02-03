package enum

import (
	"fmt"

	"github.com/gocql/gocql"
)

// MarshalCQL: Go (Int) -> Cassandra (String)
func (e StoryInteractionType) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return []byte(e.String()), nil
}

// UnmarshalCQL: Cassandra (String) -> Go (Int)
func (e *StoryInteractionType) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	val, err := StoryInteractionTypeString(string(data))
	if err != nil {
		return fmt.Errorf("unknown interaction type: %s", string(data))
	}
	*e = val
	return nil
}