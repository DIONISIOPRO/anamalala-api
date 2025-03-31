package models

import (
	"time"
)

// InformationType represents the type of information
type InformationType string

const (
	InfoTypeNews         InformationType = "news"
	InfoTypeEvent        InformationType = "event"
	InfoTypeAnnouncement InformationType = "announcement"
)

// Information represents an information post by administrators
type Information struct {
	ID          string          `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string          `bson:"title" json:"title" validate:"required"`
	Content     string          `bson:"content" json:"content" validate:"required"`
	Type        InformationType `bson:"type" json:"type"`
	AuthorID    string          `bson:"author_id" json:"author_id"`
	Attachments []string        `bson:"attachments,omitempty" json:"attachments,omitempty"`
	Published   bool            `bson:"published" json:"published"`
	CreatedAt   time.Time       `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time       `bson:"updated_at" json:"updated_at"`
	PublishedAt time.Time       `bson:"published_at,omitempty" json:"published_at,omitempty"`
}

// InformationResponse represents the information data returned to clients
type InformationResponse struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Content     string          `json:"content"`
	Type        InformationType `json:"type"`
	Author      UserResponse    `json:"author"`
	Attachments []string        `json:"attachments,omitempty"`
	Published   bool            `json:"published"`
	CreatedAt   time.Time       `json:"created_at"`
	PublishedAt time.Time       `json:"published_at,omitempty"`
}

// InformationCreation represents data for creating a new information post
type InformationCreation struct {
	Title       string          `json:"title" validate:"required"`
	Content     string          `json:"content" validate:"required"`
	Type        InformationType `json:"type" validate:"required"`
	Attachments []string        `json:"attachments,omitempty"`
	Published   bool            `json:"published"`
}

// Informations represents a slice of Information
type Informations []Information
