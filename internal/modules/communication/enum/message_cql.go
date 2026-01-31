package enum

import (
	"fmt"

	"github.com/gocql/gocql"
)

// MarshalCQL: Go (Int) -> Cassandra (String)
func (e MessageType) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return []byte(e.String()), nil
}

// UnmarshalCQL: Cassandra (String) -> Go (Int)
func (e *MessageType) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	val, err := MessageTypeString(string(data))
	if err != nil {
		return fmt.Errorf("unknown message type: %s", string(data))
	}
	*e = val
	return nil
}