package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Collections names
const (
	UsersCollection        = "users"
	PostsCollection        = "posts"
	CommentsCollection     = "comments"
	InformationCollection  = "information"
	SuggestionsCollection  = "suggestions"
	NotificationsCollection = "notifications"
)

// Client represents a MongoDB client with its database
type Client struct {
	client   *mongo.Client
	database *mongo.Database
}

// Connect establishes a connection to MongoDB
func Connect(ctx context.Context, uri, dbName string) (*Client, error) {
	clientOptions := options.Client().ApplyURI(uri)
	
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping to verify connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return &Client{
		client:   client,
		database: client.Database(dbName),
	}, nil
}

// Close closes the MongoDB connection
func (c *Client) Close(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}

// GetCollection returns a collection from the database
func (c *Client) GetCollection(name string) *mongo.Collection {
	return c.database.Collection(name)
}

// CreateIndexes creates indexes for better query performance
func (c *Client) CreateIndexes(ctx context.Context) error {
	// User indexes
	userCollection := c.GetCollection(UsersCollection)
	userIndexes := []mongo.IndexModel{
		{
			Keys: map[string]interface{}{
				"contact": 1,
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: map[string]interface{}{
				"province": 1,
			},
		},
		{
			Keys: map[string]interface{}{
				"password_reset.token": 1,
			},
		},
	}
	_, err := userCollection.Indexes().CreateMany(ctx, userIndexes)
	if err != nil {
		return err
	}

	// Post indexes
	postCollection := c.GetCollection(PostsCollection)
	postIndexes := []mongo.IndexModel{
		{
			Keys: map[string]interface{}{
				"user_id": 1,
			},
		},
		{
			Keys: map[string]interface{}{
				"created_at": -1,
			},
		},
	}
	_, err = postCollection.Indexes().CreateMany(ctx, postIndexes)
	if err != nil {
		return err
	}

	// Comment indexes
	commentCollection := c.GetCollection(CommentsCollection)
	commentIndexes := []mongo.IndexModel{
		{
			Keys: map[string]interface{}{
				"post_id": 1,
			},
		},
		{
			Keys: map[string]interface{}{
				"user_id": 1,
			},
		},
	}
	_, err = commentCollection.Indexes().CreateMany(ctx, commentIndexes)
	if err != nil {
		return err
	}

	// Information indexes
	infoCollection := c.GetCollection(InformationCollection)
	infoIndexes := []mongo.IndexModel{
		{
			Keys: map[string]interface{}{
				"published": 1,
			},
		},
		{
			Keys: map[string]interface{}{
				"author_id": 1,
			},
		},
		{
			Keys: map[string]interface{}{
				"type": 1,
			},
		},
		{
			Keys: map[string]interface{}{
				"created_at": -1,
			},
		},
	}
	_, err = infoCollection.Indexes().CreateMany(ctx, infoIndexes)
	if err != nil {
		return err
	}

	// Suggestion indexes
	suggestionCollection := c.GetCollection(SuggestionsCollection)
	suggestionIndexes := []mongo.IndexModel{
		{
			Keys: map[string]interface{}{
				"user_id": 1,
			},
		},
		{
			Keys: map[string]interface{}{
				"status": 1,
			},
		},
	}
	_, err = suggestionCollection.Indexes().CreateMany(ctx, suggestionIndexes)
	if err != nil {
		return err
	}

	// Notification indexes
	notificationCollection := c.GetCollection(NotificationsCollection)
	notificationIndexes := []mongo.IndexModel{
		{
			Keys: map[string]interface{}{
				"user_id": 1,
			},
		},
		{
			Keys: map[string]interface{}{
				"read": 1,
			},
		},
		{
			Keys: map[string]interface{}{
				"created_at": -1,
			},
		},
	}
	_, err = notificationCollection.Indexes().CreateMany(ctx, notificationIndexes)
	if err != nil {
		return err
	}

	return nil
}
