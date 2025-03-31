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

// InformationRepository implements the interfaces.InformationRepository interface
type InformationRepository struct {
	collection *mongo.Collection
}

// NewInformationRepository creates a new InformationRepository
func NewInformationRepository(client *Client) *InformationRepository {
	return &InformationRepository{
		collection: client.GetCollection(InformationCollection),
	}
}

// Create inserts a new information post into the database
func (r *InformationRepository) Create(ctx context.Context, info models.Information) error {
	info.ID = primitive.NewObjectID().Hex()
	info.CreatedAt = time.Now()
	info.UpdatedAt = time.Now()

	if info.Published {
		now := time.Now()
		info.PublishedAt = now
	}

	_, err := r.collection.InsertOne(ctx, info)
	return err
}

// FindByID finds an information post by ID
func (r *InformationRepository) FindByID(ctx context.Context, id string) (models.Information, error) {
	var info models.Information

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&info)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Information{}, nil
		}
		return models.Information{}, err
	}

	return info, nil
}

// Update updates an information post
func (r *InformationRepository) Update(ctx context.Context, info models.Information) error {
	info.UpdatedAt = time.Now()

	filter := bson.M{"_id": info.ID}
	update := bson.M{"$set": info}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete deletes an information post by ID
func (r *InformationRepository) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

// List returns a paginated list of information posts
func (r *InformationRepository) List(ctx context.Context, page, limit int64, onlyPublished bool) (models.Informations, int64, error) {
	var infos models.Informations

	filter := bson.M{}
	if onlyPublished {
		filter["published"] = true
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"created_at": -1})
	findOptions.SetSkip((page - 1) * limit)
	findOptions.SetLimit(limit)

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &infos); err != nil {
		return nil, 0, err
	}

	// Count total information posts
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return infos, total, nil
}

// Publish publishes an information post
func (r *InformationRepository) Publish(ctx context.Context, id string) error {
	now := time.Now()

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"published":    true,
			"published_at": now,
			"updated_at":   now,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Unpublish unpublishes an information post
func (r *InformationRepository) Unpublish(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"published":  false,
			"updated_at": time.Now(),
		},
		"$unset": bson.M{"published_at": ""},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
