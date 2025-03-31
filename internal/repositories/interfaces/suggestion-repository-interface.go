package interfaces

import (
	"context"

	"github.com/anamalala/internal/models"
)

// SuggestionRepository defines the interface for suggestion repository
type SuggestionRepository interface {
	Create(ctx context.Context, suggestion models.Suggestion) error
	FindByID(ctx context.Context, id string) (models.Suggestion, error)
	Update(ctx context.Context, suggestion models.Suggestion) error
	Delete(ctx context.Context, id string) error
	GetByStatus(ctx context.Context, status string, page, limit int64) (models.Suggestions, int64, error)
	List(ctx context.Context, page, limit int64, status models.SuggestionStatus) (models.Suggestions, int64, error)
	ListByUserID(ctx context.Context, userID string, page, limit int64) (models.Suggestions, int64, error)
}
