package sharedEnums

import (
	"fmt"

	"github.com/gocql/gocql"
)

// MarshalCQL: Go (Int) -> Cassandra (String)
func (e UserBadge) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return []byte(e.String()), nil
}

// UnmarshalCQL: Cassandra (String) -> Go (Int)
func (e *UserBadge) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	val, err := UserBadgeString(string(data))
	if err != nil {
		return fmt.Errorf("unknown badge: %s", string(data))
	}
	*e = val
	return nil
}
