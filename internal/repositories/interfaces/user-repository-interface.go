package interfaces

import (
	"context"
	"time"

	"github.com/anamalala/internal/models"
)

// UserRepository defines the interface for user repository
type UserRepository interface {
	Create(ctx context.Context, user models.User) error
	FindByID(ctx context.Context, id string) (models.User, error)
	FindByContact(ctx context.Context, contact string) (models.User, error)
	Update(ctx context.Context, user models.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int64) (models.Users, int64, error)
	InactiveUsers(ctx context.Context, page, limit int64) (models.Users, int64, error)
	ListByProvince(ctx context.Context, province string, page, limit int64) (models.Users, int64, error)
	UpdatePassword(ctx context.Context, id string, password string) error
	StorePasswordResetToken(ctx context.Context, contact, token string, expiryTime time.Time) error
	ValidatePasswordResetToken(ctx context.Context, token string) (models.User, error)
	ToggleUserActive(ctx context.Context, id string, active bool) error
	GetAllContacts(ctx context.Context) ([]string, error)
	GetContactsByProvince(ctx context.Context, province string) ([]string, error)
	ListByRole(ctx context.Context, role string, page, limit int64) (models.Users, int64, error)
}
