package sharedEnums

import "go.mongodb.org/mongo-driver/bson/bsontype"

func (e Type_Block) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

func (e *Type_Block) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	val, err := Type_BlockString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}
