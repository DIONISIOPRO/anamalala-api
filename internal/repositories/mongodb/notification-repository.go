package mongodb

import (
	"context"

	"github.com/anamalala/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// NotificationRepository implements the interfaces.NotificationRepository interface
type NotificationRepository struct {
	collection *mongo.Collection
}

// NewSuggestionRepository creates a new SuggestionRepository
func NewNotificationRepository(client *Client) *NotificationRepository {
	return &NotificationRepository{
		collection: client.GetCollection("notifications"),
	}
}

func (r *NotificationRepository) Create(ctx context.Context, notification models.Notification) error {
	return nil
}
func (r *NotificationRepository) CreateMany(ctx context.Context, notifications []models.Notification) error {
	return nil
}
func (r *NotificationRepository) FindByID(ctx context.Context, id string) (models.Notification, error) {
	return models.Notification{}, nil
}
func (r *NotificationRepository) MarkAsRead(ctx context.Context, id string) error {

	return nil
}
func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID string) error {

	return nil
}
func (r *NotificationRepository) Delete(ctx context.Context, id string) error {
	return nil
}
func (r *NotificationRepository) ListByUserID(ctx context.Context, userID string, page, limit int64, unreadOnly bool) (models.Notifications, int64, error) {
	return models.Notifications{}, 0, nil
}
func (r *NotificationRepository) CountUnread(ctx context.Context, userID string) (int64, error) {
	return 0, nil
}
