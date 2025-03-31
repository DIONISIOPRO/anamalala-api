package interfaces

import (
	"context"

	"github.com/anamalala/internal/models"
)

// InformationRepository defines the interface for information repository
type InformationRepository interface {
	Create(ctx context.Context, info models.Information) error
	FindByID(ctx context.Context, id string) (models.Information, error)
	Update(ctx context.Context, info models.Information) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int64, onlyPublished bool) (models.Informations, int64, error)
	Publish(ctx context.Context, id string) error
	Unpublish(ctx context.Context, id string) error
}
