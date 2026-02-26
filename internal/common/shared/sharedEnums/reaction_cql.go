package sharedEnums

import (
	"fmt"

	"github.com/gocql/gocql"
)

// =============================================================================
// HELPER: REACTION TARGET
// =============================================================================

// MarshalCQL: Go (Int) -> Cassandra (String)
func (e ReactionTarget) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return []byte(e.String()), nil
}

// UnmarshalCQL: Cassandra (String) -> Go (Int)
func (e *ReactionTarget) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	val, err := ReactionTargetString(string(data))
	if err != nil {
		return fmt.Errorf("unknown reaction target: %s", string(data))
	}
	*e = val
	return nil
}

// =============================================================================
// HELPER: REACTION CODE
// =============================================================================

// MarshalCQL: Go (Int) -> Cassandra (String)
func (e ReactionCode) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return []byte(e.String()), nil
}

// UnmarshalCQL: Cassandra (String) -> Go (Int)
func (e *ReactionCode) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	val, err := ReactionCodeString(string(data))
	if err != nil {
		return fmt.Errorf("unknown reaction code: %s", string(data))
	}
	*e = val
	return nil
}