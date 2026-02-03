package enum

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)



// -----------------------------------------------------------------------------
// SearchTargetType
// -----------------------------------------------------------------------------
func (e SearchTargetType) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil 
}

func (e *SearchTargetType) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for SearchTargetType, got %v", t)
	}
	val, err := SearchTargetTypeString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}