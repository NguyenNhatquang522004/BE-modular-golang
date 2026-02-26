package enum

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)



// -----------------------------------------------------------------------------
// RemixType
// -----------------------------------------------------------------------------
func (e RemixType) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

func (e *RemixType) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for RemixType, got %v", t)
	}
	val, err := RemixTypeString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}

