package enum

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)


// -----------------------------------------------------------------------------
// MemberRole
// -----------------------------------------------------------------------------
func (e MemberRole) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

func (e *MemberRole) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for MemberRole, got %v", t)
	}
	val, err := MemberRoleString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}

// -----------------------------------------------------------------------------
// MemberStatus
// -----------------------------------------------------------------------------
func (e MemberStatus) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

func (e *MemberStatus) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for MemberStatus, got %v", t)
	}
	val, err := MemberStatusString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}

// -----------------------------------------------------------------------------
// MemberBadge
// -----------------------------------------------------------------------------
func (e MemberBadge) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

func (e *MemberBadge) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for MemberBadge, got %v", t)
	}
	val, err := MemberBadgeString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}