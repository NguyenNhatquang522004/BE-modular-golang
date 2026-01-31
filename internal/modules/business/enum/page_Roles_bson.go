package enum

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// -----------------------------------------------------------------------------
// PageRole
// -----------------------------------------------------------------------------
func (e PageRole) MarshalBSONValue() (bsontype.Type, []byte, error) {
	// Nếu có common_bson.go: return marshalEnum(e.String())
	return bsontype.String, []byte(e.String()), nil
}

func (e *PageRole) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for PageRole, got %v", t)
	}
	val, err := PageRoleString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}