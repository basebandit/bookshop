package model

// RecordID defines a record id. Together with RecordType identifies unique records
// across all types.
type RecordID string

// RecordType defines a record type. Togethe with RecordID identifies unique records
// across all types.
type RecordType string

// RecordTypeBook denotes existing record types.
const RecordTypeBook = RecordType("book")

// UserID defines a user id.
type UserID string

// RatingValue defines a value of a rating record.
type RatingValue int

// Rating defines an individual rating created by a user for some record.
type Rating struct {
	RecordID   string      `json:"recordID"`
	RecordType string      `json:"recordType"`
	UserID     UserID      `json:"userId"`
	Value      RatingValue `json:"value"`
}
