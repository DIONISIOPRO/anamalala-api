package interfaces

import (
	"context"

	"github.com/anamalala/internal/models"
)

// PostRepository defines the interface for post repository
type PostRepository interface {
	Create(ctx context.Context, post models.Post) (models.Post, error)
	FindByID(ctx context.Context, id string) (models.Post, error)
	Update(ctx context.Context, post models.Post) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int64) (models.Posts, int64, error)
	AddComment(ctx context.Context, postID, commentID string) error
	RemoveComment(ctx context.Context, postID, commentID string) error
	AddLike(ctx context.Context, postID, userId string) error
	RemoveLike(ctx context.Context, postID, userId string) error
}
