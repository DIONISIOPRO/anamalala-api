package models

import (
	"time"
)

// SuggestionStatus represents the status of a suggestion
type SuggestionStatus string

const (
	SuggestionStatusAll      SuggestionStatus = "all"
	SuggestionStatusNew      SuggestionStatus = "new"
	SuggestionStatusReviewed SuggestionStatus = "reviewed"
	SuggestionStatusApproved SuggestionStatus = "approved"
	SuggestionStatusRejected SuggestionStatus = "rejected"
)

// Suggestion represents a suggestion for improving the platform
type Suggestion struct {
	ID          string           `bson:"_id,omitempty" json:"id,omitempty"`
	UserID      string           `bson:"user_id" json:"user_id"`
	Title       string           `bson:"title" json:"title" validate:"required"`
	Description string           `bson:"description" json:"description" validate:"required"`
	Status      SuggestionStatus `bson:"status" json:"status"`
	AdminNotes  string           `bson:"admin_notes,omitempty" json:"admin_notes,omitempty"`
	ReviewedBy  string           `bson:"reviewed_by,omitempty" json:"reviewed_by,omitempty"`
	CreatedAt   time.Time        `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time        `bson:"updated_at" json:"updated_at"`
	ReviewedAt  time.Time        `bson:"reviewed_at,omitempty" json:"reviewed_at,omitempty"`
}

// SuggestionResponse represents the suggestion data returned to clients
type SuggestionResponse struct {
	ID          string           `json:"id"`
	User        UserResponse     `json:"user"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Status      SuggestionStatus `json:"status"`
	CreatedAt   time.Time        `json:"created_at"`
	ReviewedAt  time.Time        `json:"reviewed_at,omitempty"`
}

// SuggestionCreation represents data for creating a new suggestion
type SuggestionCreation struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

// SuggestionUpdate represents data for updating a suggestion status by admin
type SuggestionUpdate struct {
	Status     SuggestionStatus `json:"status" validate:"required"`
	AdminNotes string           `json:"admin_notes"`
}

// Suggestions represents a slice of Suggestion
type Suggestions []Suggestion
