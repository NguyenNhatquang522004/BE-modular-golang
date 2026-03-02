package enum

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// Helper (Nếu chưa có common_bson.go)

// -----------------------------------------------------------------------------
// PageStatus
// -----------------------------------------------------------------------------
func (e PageStatus) MarshalBSONValue() (bsontype.Type, []byte, error) { return bsontype.String, []byte(e.String()), nil }
func (e *PageStatus) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String { return fmt.Errorf("expected string, got %v", t) }
	val, err := PageStatusString(string(data)); if err != nil { return err }; *e = val; return nil
}

// -----------------------------------------------------------------------------
// CTAType
// -----------------------------------------------------------------------------
func (e CTAType) MarshalBSONValue() (bsontype.Type, []byte, error) { return bsontype.String, []byte(e.String()), nil }
func (e *CTAType) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String { return fmt.Errorf("expected string, got %v", t) }
	val, err := CTATypeString(string(data)); if err != nil { return err }; *e = val; return nil
}

// -----------------------------------------------------------------------------
// MessagingStatus
// -----------------------------------------------------------------------------
func (e MessagingStatus) MarshalBSONValue() (bsontype.Type, []byte, error) { return bsontype.String, []byte(e.String()), nil }
func (e *MessagingStatus) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String { return fmt.Errorf("expected string, got %v", t) }
	val, err := MessagingStatusString(string(data)); if err != nil { return err }; *e = val; return nil
}

// -----------------------------------------------------------------------------
// DayOfWeek
// -----------------------------------------------------------------------------
func (e DayOfWeek) MarshalBSONValue() (bsontype.Type, []byte, error) { return bsontype.String, []byte(e.String()), nil }
func (e *DayOfWeek) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String { return fmt.Errorf("expected string, got %v", t) }
	val, err := DayOfWeekString(string(data)); if err != nil { return err }; *e = val; return nil
}