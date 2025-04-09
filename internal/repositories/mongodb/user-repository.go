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

// UserRepository implements the interfaces.UserRepository interface
type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(client *Client) *UserRepository {
	return &UserRepository{
		collection: client.GetCollection(UsersCollection),
	}
}

// Create inserts a new user into the database
func (r *UserRepository) Create(ctx context.Context, user models.User) error {
	user.ID = primitive.NewObjectID().Hex()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	if user.Role == "" {
		user.Role = models.RoleUser
	}
	user.Active = true
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"contact": user.Contact}
	update := bson.M{"$set": user}
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id string) (models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, nil
		}
		return models.User{}, err
	}

	return user, nil
}

// FindByContact finds a user by contact (phone number)
func (r *UserRepository) FindByContact(ctx context.Context, contact string) (models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"contact": contact}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, err
		}

		return models.User{}, err

	}
	return user, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user models.User) error {
	user.UpdatedAt = time.Now()
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

// List returns a paginated list of users
func (r *UserRepository) List(ctx context.Context, page, limit int64) (models.Users, int64, error) {
	var users models.Users

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"created_at": -1})
	findOptions.SetSkip((page - 1) * limit)
	findOptions.SetLimit(limit)

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	// Count total users
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// ListByProvince returns a paginated list of users by province
func (r *UserRepository) ListByProvince(ctx context.Context, province string, page, limit int64) (models.Users, int64, error) {
	var users models.Users

	filter := bson.M{"province": province}

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"created_at": -1})
	findOptions.SetSkip((page - 1) * limit)
	findOptions.SetLimit(limit)

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	// Count total users by province
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) ListByRole(ctx context.Context, role string, page, limit int64) (models.Users, int64, error) {
	var users models.Users

	filter := bson.M{"role": role}

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"created_at": -1})
	findOptions.SetSkip((page - 1) * limit)
	findOptions.SetLimit(limit)

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	// Count total users by province
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(ctx context.Context, id string, password string) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"password":   password,
			"updated_at": time.Now(),
		},
		"$unset": bson.M{"password_reset": ""},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// StorePasswordResetToken stores a password reset token for a user
func (r *UserRepository) StorePasswordResetToken(ctx context.Context, contact, token string, expiryTime time.Time) error {
	filter := bson.M{"contact": contact}
	update := bson.M{
		"$set": bson.M{
			"password_reset.token":      token,
			"password_reset.expires_at": expiryTime,
			"updated_at":                time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no user found with the provided contact")
	}

	return nil
}

// ValidatePasswordResetToken validates a password reset token
func (r *UserRepository) ValidatePasswordResetToken(ctx context.Context, token string) (models.User, error) {
	var user models.User
	filter := bson.M{
		"password_reset.token":      token,
		"password_reset.expires_at": bson.M{"$gt": primitive.NewDateTimeFromTime(time.Now())},
	}

	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, errors.New("invalid or expired token")
		}
		return models.User{}, err
	}
	return user, nil
}

// ToggleUserActive activates or deactivates a user
func (r *UserRepository) ToggleUserActive(ctx context.Context, id string, active bool) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"active":     active,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// GetAllContacts returns all user contact numbers
func (r *UserRepository) GetAllContacts(ctx context.Context) ([]string, error) {
	var contacts []string
	projection := bson.M{"contact": 1, "_id": 0}
	findOptions := options.Find().SetProjection(projection)
	cursor, err := r.collection.Find(ctx, bson.M{"active": true}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	type contactDoc struct {
		Contact string `bson:"contact"`
	}

	var results []contactDoc
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	contacts = make([]string, len(results))
	for i, doc := range results {
		contacts[i] = doc.Contact
	}

	return contacts, nil
}

// GetContactsByProvince returns all user contact numbers for a province
func (r *UserRepository) GetContactsByProvince(ctx context.Context, province string) ([]string, error) {
	var contacts []string
	filter := bson.M{"province": province, "active": true}
	projection := bson.M{"contact": 1, "_id": 0}
	findOptions := options.Find().SetProjection(projection)
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	type contactDoc struct {
		Contact string `bson:"contact"`
	}

	var results []contactDoc
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	contacts = make([]string, len(results))
	for i, doc := range results {
		contacts[i] = doc.Contact
	}

	return contacts, nil
}

func (r *UserRepository) InactiveUsers(ctx context.Context, page, limit int64) (models.Users, int64, error) {
	var users models.Users

	filter := bson.M{"active": false}

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"created_at": -1})
	findOptions.SetSkip((page - 1) * limit)
	findOptions.SetLimit(limit)

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	// Count total users by province
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
