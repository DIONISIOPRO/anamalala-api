package interfaces

import (
	"context"

	"github.com/anamalala/internal/models"
)

// CommentRepository defines the interface for comment repository
type CommentRepository interface {
	Create(ctx context.Context, comment models.Comment) error
	FindByID(ctx context.Context, id string) (models.Comment, error)
	Update(ctx context.Context, comment models.Comment) error
	RemoveLike(ctx context.Context, commentObjectID, userObjectID string) error
	AddLike(ctx context.Context, commentObjectID, userObjectID string) error
	Delete(ctx context.Context, id string) error
	ListByPostID(ctx context.Context, postID string, page, limit int64) (models.Comments, int64, error)
	ListByCommentID(ctx context.Context, postID string, page, limit int64) (models.Comments, int64, error)
}
