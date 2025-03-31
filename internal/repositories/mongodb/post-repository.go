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

// PostRepository implements the interfaces.PostRepository interface
type PostRepository struct {
	collection *mongo.Collection
}

// NewPostRepository creates a new PostRepository
func NewPostRepository(client *Client) *PostRepository {
	return &PostRepository{
		collection: client.GetCollection(PostsCollection),
	}
}

// Create inserts a new post into the database
func (r *PostRepository) Create(ctx context.Context, post models.Post) (models.Post, error) {
	post.ID = primitive.NewObjectID().Hex()
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	post.Comments = []models.Comment{}
	post.LikedUserId = []string{}
	post.Likes = 0
	_, err := r.collection.InsertOne(ctx, post)

	return post, err
}

// FindByID finds a post by ID
func (r *PostRepository) FindByID(ctx context.Context, id string) (models.Post, error) {
	var post models.Post

	filter := bson.M{
		"_id":        id,
		"deleted_at": nil,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&post)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Post{}, nil
		}
		return models.Post{}, err
	}

	return post, nil
}

// Update updates a post
func (r *PostRepository) Update(ctx context.Context, post models.Post) error {
	post.UpdatedAt = time.Now()

	filter := bson.M{"_id": post.ID}
	update := bson.M{"$set": post}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete soft deletes a post by ID
func (r *PostRepository) Delete(ctx context.Context, id string) error {
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

// List returns a paginated list of posts
func (r *PostRepository) List(ctx context.Context, page, limit int64) (models.Posts, int64, error) {

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"created_at": -1})
	findOptions.SetSkip((page - 1) * limit)
	findOptions.SetLimit(limit)
    // Execute aggregation
    cursor, err := r.collection.Find(ctx, bson.D{},findOptions )
    if err != nil {
        return nil, 0, err
    }
    defer cursor.Close(ctx)

    // Decode results
    var posts []models.Post
    if err = cursor.All(ctx, &posts); err != nil {
        return nil, 0, err
    }

    // Count total posts
    totalPosts, err := r.collection.CountDocuments(ctx, bson.D{})
    if err != nil {
        return nil, 0, err
    }

	
    return posts, totalPosts, nil
}





// AddComment adds a comment ID to a post's comments array
func (r *PostRepository) AddComment(ctx context.Context, postID, commentID string) error {
	filter := bson.M{"_id": postID}
	update := bson.M{
		"$push": bson.M{"comments": commentID},
		"$set":  bson.M{"updated_at": time.Now()},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// RemoveComment removes a comment ID from a post's comments array
func (r *PostRepository) RemoveComment(ctx context.Context, postID, commentID string) error {
	filter := bson.M{"_id": postID}
	update := bson.M{
		"$pull": bson.M{"comments": commentID},
		"$set":  bson.M{"updated_at": time.Now()},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *PostRepository) AddLike(ctx context.Context, postID, userId string) error {
	post, err := r.FindByID(ctx, postID)
	if err != nil {
		return err
	}
	post.Likes = post.Likes + 1
	post.LikedUserId = append(post.LikedUserId, string(userId))

	err = r.Update(ctx, post)
	if err != nil {
		return err
	}
	return err
}

func (r *PostRepository) RemoveLike(ctx context.Context, postID, userId string) error {
	post, err := r.FindByID(ctx, postID)
	if err != nil {
		return err
	}
	if post.Likes > 0 {
		post.Likes = post.Likes - 1
	}
	for i := 0; i < len(post.LikedUserId); {
		if post.LikedUserId[i] == userId {
			if len(post.LikedUserId) > i {
				post.LikedUserId = append(post.LikedUserId[:i], post.LikedUserId[i+1:]...)
				continue
			}
		}
		i++
	}

	err = r.Update(ctx, post)
	if err != nil {
		return err
	}
	return err
}
