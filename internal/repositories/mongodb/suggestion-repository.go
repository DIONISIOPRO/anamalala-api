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

// SuggestionRepository implements the interfaces.SuggestionRepository interface
type SuggestionRepository struct {
	collection *mongo.Collection
}

// NewSuggestionRepository creates a new SuggestionRepository
func NewSuggestionRepository(client *Client) *SuggestionRepository {
	return &SuggestionRepository{
		collection: client.GetCollection(SuggestionsCollection),
	}
}

// Create inserts a new suggestion into the database
func (r *SuggestionRepository) Create(ctx context.Context, suggestion models.Suggestion) error {
	suggestion.ID = primitive.NewObjectID().Hex()
	suggestion.CreatedAt = time.Now()
	suggestion.UpdatedAt = time.Now()
	
	if suggestion.Status == "" {
		suggestion.Status = models.SuggestionStatusNew
	}
	
	_, err := r.collection.InsertOne(ctx, suggestion)
	return err
}

// FindByID finds a suggestion by ID
func (r *SuggestionRepository) FindByID(ctx context.Context, id string) (models.Suggestion, error) {
	var suggestion models.Suggestion
	
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&suggestion)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Suggestion{}, nil
		}
		return models.Suggestion{}, err
	}
	
	return suggestion, nil
}

// Update updates a suggestion
func (r *SuggestionRepository) Update(ctx context.Context, suggestion models.Suggestion) error {
	suggestion.UpdatedAt = time.Now()
	
	if suggestion.Status != models.SuggestionStatusNew {
		now := time.Now()
		suggestion.ReviewedAt = now
	}
	
	filter := bson.M{"_id": suggestion.ID}
	update := bson.M{"$set": suggestion}
	
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// List returns a paginated list of suggestions, optionally filtered by status
func (r *SuggestionRepository) List(ctx context.Context, page, limit int64, status models.SuggestionStatus) (models.Suggestions, int64, error) {
	var suggestions models.Suggestions
	
	filter := bson.M{}
	if status != "" {
		filter["status"] = status
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
	
	if err := cursor.All(ctx, &suggestions); err != nil {
		return nil, 0, err
	}
	
	// Count total suggestions
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	
	return suggestions, total, nil
}

// ListByUserID returns a paginated list of suggestions for a user
func (r *SuggestionRepository) ListByUserID(ctx context.Context, userID string, page, limit int64) (models.Suggestions, int64, error) {
	var suggestions models.Suggestions
	
	filter := bson.M{"user_id": userID}
	
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"created_at": -1})
	findOptions.SetSkip((page - 1) * limit)
	findOptions.SetLimit(limit)
	
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	
	if err := cursor.All(ctx, &suggestions); err != nil {
		return nil, 0, err
	}
	
	// Count total suggestions for a user
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	
	return suggestions, total, nil
}
func (r *SuggestionRepository) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}


func (r *SuggestionRepository) GetByStatus(ctx context.Context, status string,  page, limit int64)  (models.Suggestions, int64, error){
	var suggestions models.Suggestions
	
	filter := bson.M{"status": status}
	
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"created_at": -1})
	findOptions.SetSkip((page - 1) * limit)
	findOptions.SetLimit(limit)
	
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	
	if err := cursor.All(ctx, &suggestions); err != nil {
		return nil, 0, err
	}
	
	// Count total suggestions for a user
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	
	return suggestions, total, nil
}