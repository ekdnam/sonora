package psql

import (
	"database/sql"
	"fmt"
	"time"
)

// CreateTopicValidation inserts a new validation record
func (db *DB) CreateTopicValidation(topic string, isValid bool, reason string) (*TopicValidation, error) {
	validation := &TopicValidation{
		Topic:     topic,
		IsValid:   isValid,
		Reason:    reason,
		CreatedAt: time.Now(),
	}

	query := `
		INSERT INTO topic_validations (topic, is_valid, reason, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	err := db.conn.QueryRow(query, validation.Topic, validation.IsValid,
		validation.Reason, validation.CreatedAt).Scan(&validation.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create topic validation: %w", err)
	}
	return validation, nil
}

// GetTopicValidationByID retrieves a topic validation record by its ID
func (db *DB) GetTopicValidationByID(id int64) (*TopicValidation, error) {
	validation := &TopicValidation{}

	query := `
		SELECT id, topic, is_valid, reason, created_at
		FROM topic_validations
		WHERE id = $1`

	err := db.conn.QueryRow(query, id).Scan(
		&validation.ID,
		&validation.Topic,
		&validation.IsValid,
		&validation.Reason,
		&validation.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("topic validation with id %d not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get topic validation: %w", err)
	}

	return validation, nil
}

// GetTopicValidationByTopic retrieves the most recent topic validation record for a given topic
func (db *DB) GetTopicValidationByTopic(topic string) (*TopicValidation, error) {
	validation := &TopicValidation{}

	query := `
		SELECT id, topic, is_valid, reason, created_at
		FROM topic_validations
		WHERE topic = $1
		ORDER BY created_at DESC
		LIMIT 1`

	err := db.conn.QueryRow(query, topic).Scan(
		&validation.ID,
		&validation.Topic,
		&validation.IsValid,
		&validation.Reason,
		&validation.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no validation found for topic: %s", topic)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get topic validation: %w", err)
	}

	return validation, nil
}
