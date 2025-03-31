package interfaces

import (
	"context"

	"github.com/anamalala/internal/models"
)

// NotificationRepository defines the interface for notification repository
type NotificationRepository interface {
	Create(ctx context.Context, notification models.Notification) error
	CreateMany(ctx context.Context, notifications []models.Notification) error
	FindByID(ctx context.Context, id string) (models.Notification, error)
	MarkAsRead(ctx context.Context, id string) error
	MarkAllAsRead(ctx context.Context, userID string) error
	Delete(ctx context.Context, id string) error
	ListByUserID(ctx context.Context, userID string, page, limit int64, unreadOnly bool) (models.Notifications, int64, error)
	CountUnread(ctx context.Context, userID string) (int64, error)
}
