package models

import (
	"time"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeSystem NotificationType = "system"
	NotificationTypeChat   NotificationType = "chat"
	NotificationTypeInfo   NotificationType = "info"
	NotificationTypeAdmin  NotificationType = "admin"
)

// Notification represents a notification for a user
type Notification struct {
	ID        string           `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    string           `bson:"user_id" json:"user_id"`
	Type      NotificationType `bson:"type" json:"type"`
	Title     string           `bson:"title" json:"title"`
	Message   string           `bson:"message" json:"message"`
	Read      bool             `bson:"read" json:"read"`
	Reference string           `bson:"reference,omitempty" json:"reference,omitempty"`
	CreatedAt time.Time        `bson:"created_at" json:"created_at"`
	ReadAt    time.Time       `bson:"read_at,omitempty" json:"read_at,omitempty"`
}

// NotificationResponse represents the notification data returned to clients
type NotificationResponse struct {
	ID        string           `json:"id"`
	Type      NotificationType `json:"type"`
	Title     string           `json:"title"`
	Message   string           `json:"message"`
	Read      bool             `json:"read"`
	Reference string          `json:"reference,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
	ReadAt    time.Time       `json:"read_at,omitempty"`
}

// NotificationCreation represents data for creating a new notification
type NotificationCreation struct {
	UserID    string           `json:"user_id" validate:"required"`
	Type      NotificationType `json:"type" validate:"required"`
	Title     string           `json:"title" validate:"required"`
	Message   string           `json:"message" validate:"required"`
	Reference string          `json:"reference,omitempty"`
}

// BulkNotificationCreation represents data for creating bulk notifications
type BulkNotificationCreation struct {
	UserIDs   []string         `json:"user_ids" validate:"required"`
	Type      NotificationType `json:"type" validate:"required"`
	Title     string           `json:"title" validate:"required"`
	Message   string           `json:"message" validate:"required"`
	Reference string          `json:"reference,omitempty"`
}

// ProvinceNotificationCreation represents data for creating notifications for users in a province
type ProvinceNotificationCreation struct {
	Province  string           `json:"province" validate:"required"`
	Type      NotificationType `json:"type" validate:"required"`
	Title     string           `json:"title" validate:"required"`
	Message   string           `json:"message" validate:"required"`
	Reference string          `json:"reference,omitempty"`
}

// Notifications represents a slice of Notification
type Notifications []Notification
