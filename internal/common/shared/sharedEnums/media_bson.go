package sharedEnums

// import (
// 	"fmt"

// 	"go.mongodb.org/mongo-driver/bson/bsontype"
// )

// // MarshalBSONValue: Int -> String
// func (e MediaType) MarshalBSONValue() (bsontype.Type, []byte, error) {
// 	return bsontype.String, []byte(e.String()), nil
// }

// // UnmarshalBSONValue: String -> Int
// func (e *MediaType) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
// 	if t != bsontype.String {
// 		return fmt.Errorf("expected string for MediaType, got %v", t)
// 	}
// 	// MediaTypeString được sinh ra bởi lệnh 'go generate'
// 	val, err := MediaTypeString(string(data))
// 	if err != nil {
// 		return err
// 	}
// 	*e = val
// 	return nil
// }
