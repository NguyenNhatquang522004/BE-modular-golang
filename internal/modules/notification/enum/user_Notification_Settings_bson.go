package enum

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)


// -----------------------------------------------------------------------------
// EmailFrequency
// -----------------------------------------------------------------------------
func (e EmailFrequency) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil 
}

func (e *EmailFrequency) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for EmailFrequency, got %v", t)
	}
	val, err := EmailFrequencyString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}