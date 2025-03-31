package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/anamalala/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CommentRepository implements the interfaces.CommentRepository interface
type CommentRepository struct {
	collection *mongo.Collection
}

// NewCommentRepository creates a new CommentRepository
func NewCommentRepository(client *Client) *CommentRepository {
	return &CommentRepository{
		collection: client.GetCollection(CommentsCollection),
	}
}

// Create inserts a new comment into the database
func (r *CommentRepository) Create(ctx context.Context, comment models.Comment) error {
	comment.ID = primitive.NewObjectID().Hex()
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, comment)
	return err
}

// FindByID finds a comment by ID
func (r *CommentRepository) FindByID(ctx context.Context, id string) (models.Comment, error) {
	var comment models.Comment
	filter := bson.M{
		"_id":        id,
		"deleted_at": nil,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&comment)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Comment{}, nil
		}
		return models.Comment{}, err
	}

	return comment, nil
}

// Update updates a comment
func (r *CommentRepository) Update(ctx context.Context, comment models.Comment) error {
	comment.UpdatedAt = time.Now()

	filter := bson.M{"_id": comment.ID}
	update := bson.M{"$set": comment}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete soft deletes a comment by ID
func (r *CommentRepository) Delete(ctx context.Context, id string) error {
	now := time.Now()

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"deleted_at": now,
			"updated_at": now,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// ListByPostID returns a paginated list of comments for a post
func (r *CommentRepository) ListByCommentID(ctx context.Context, postID string, page, limit int64) (models.Comments, int64, error) {
	var comments models.Comments

	filter := bson.M{
		"reference_id":    postID,
		"reference":    "comment",
		"deleted_at": nil,
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"created_at": 1}) // Oldest first
	findOptions.SetSkip((page - 1) * limit)
	findOptions.SetLimit(limit)

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &comments); err != nil {
		return nil, 0, err
	}

	// Count total active comments for a post
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// ListByPostID returns a paginated list of comments for a post
func (r *CommentRepository) ListByPostID(ctx context.Context, postID string, page, limit int64) (models.Comments, int64, error) {
	var comments models.Comments

	filter := bson.M{
		"reference_id":    postID,
		"reference":    "post",
		"deleted_at": nil,
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"created_at": 1}) // Oldest first
	findOptions.SetSkip((page - 1) * limit)
	findOptions.SetLimit(limit)

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &comments); err != nil {
		return nil, 0, err
	}

	// Count total active comments for a post
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

func (r *CommentRepository) RemoveLike(ctx context.Context, commentObjectID, userObjectID string) error {
	now := time.Now()
	comment := models.Comment{}
	filter := bson.M{"_id": commentObjectID}
	result := r.collection.FindOne(ctx, filter)
	err := result.Decode(comment)
	if err != nil {
		return err
	}
	if comment.Likes > 0 {
		comment.Likes = comment.Likes - 1
	}
	comment.UpdatedAt = now
	_, err = r.collection.UpdateOne(ctx, filter, comment)
	return err
}

func (r *CommentRepository) AddLike(ctx context.Context, commentObjectID, userObjectID string) error {
	now := time.Now()

	filter := bson.M{"_id": commentObjectID}
	update := bson.M{
		"$inc": bson.M{
			"likes":      1,
			"updated_at": now,
		},

		"$push": bson.M{
			"likeduserid": userObjectID,
			"updated_at":  now,
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
