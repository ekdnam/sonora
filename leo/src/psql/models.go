package psql

import (
	"time"
)

// TopicValidation represents a record of topic validation attempts
// Table name: topic_validations
type TopicValidation struct {
	ID        int64     `db:"id"`         // Primary key
	Topic     string    `db:"topic"`      // Input topic that was validated
	IsValid   bool      `db:"is_valid"`   // Validation response
	Reason    string    `db:"reason"`     // Explanation for the validation result
	CreatedAt time.Time `db:"created_at"` // Timestamp of validation
}

// TableName returns the name of the database table
func (TopicValidation) TableName() string {
	return "topic_validations"
}
