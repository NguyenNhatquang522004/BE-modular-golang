package sharedEnums

import (
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// Nếu bạn có dùng MongoDB (như trong file mẫu có BSON) thì thêm 2 hàm này,
// nếu chỉ dùng Postgres thuần thì có thể bỏ qua 2 hàm Marshal/Unmarshal BSON này.
func (e StatusFriendship) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

func (e *StatusFriendship) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	val, err := StatusFriendshipString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}
