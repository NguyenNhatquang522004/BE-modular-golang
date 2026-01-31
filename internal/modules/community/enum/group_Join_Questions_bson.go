package enum

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// -----------------------------------------------------------------------------
// QuestionType
// -----------------------------------------------------------------------------
func (e QuestionType) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

func (e *QuestionType) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for QuestionType, got %v", t)
	}
	val, err := QuestionTypeString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}