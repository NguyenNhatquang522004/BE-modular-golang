package sharedEnums

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

func (e MusicGenre) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

func (e *MusicGenre) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for MusicGenre, got %v", t)
	}
	val, err := MusicGenreString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}
