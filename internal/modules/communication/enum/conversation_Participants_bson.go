package enum

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// -----------------------------------------------------------------------------
// ParticipantRole
// -----------------------------------------------------------------------------
func (e ParticipantRole) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

func (e *ParticipantRole) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for ParticipantRole, got %v", t)
	}
	val, err := ParticipantRoleString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}