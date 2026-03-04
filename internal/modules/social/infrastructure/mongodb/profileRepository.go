package mongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProfileRepository struct {
	client *mongo.Database
}

func NewProfileRepository(client *mongo.Database) *ProfileRepository {
	return &ProfileRepository{
		client: client,
	}
}

func (r *ProfileRepository) CreateProfile(ctx context.Context, profileData *entity.Profiles) error {
	// Implement logic to create a new profile in MongoDB

	collection := r.client.Collection(entity.Profiles{}.CollectionNameProfiles())
	_, err := collection.InsertOne(ctx, profileData)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProfileRepository) GetProfileByID(ctx context.Context, profileID string) (*entity.Profiles, error) {
	// Implement logic to retrieve a profile by ID from MongoDB

	collection := r.client.Collection(entity.Profiles{}.CollectionNameProfiles())
	filter := bson.M{"_id": profileID}
	var profile entity.Profiles
	err := collection.FindOne(ctx, filter).Decode(&profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *ProfileRepository) GetProfileByUserID(ctx context.Context, userID string) (*entity.Profiles, error) {
	// Implement logic to retrieve a profile by UserID from MongoDB
	collection := r.client.Collection(entity.Profiles{}.CollectionNameProfiles())
	filter := bson.M{"user_id": userID}
	var profile entity.Profiles
	err := collection.FindOne(ctx, filter).Decode(&profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *ProfileRepository) UpdateProfile(ctx context.Context, profileData *entity.Profiles) error {
	collection := r.client.Collection(entity.Profiles{}.CollectionNameProfiles())
	filter := bson.M{"_id": profileData.ID}
	err := collection.FindOneAndReplace(ctx, filter, profileData).Err()
	if err != nil {
		return err
	}
	return nil
}
func (r *ProfileRepository) DeleteProfile(ctx context.Context, profileID string) error {
	collection := r.client.Collection(entity.Profiles{}.CollectionNameProfiles())
	filter := bson.M{"_id": profileID}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
