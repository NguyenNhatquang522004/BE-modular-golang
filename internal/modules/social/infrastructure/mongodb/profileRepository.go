package mongodb

import (
	"context"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/delivery/mapper"
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

func (r *ProfileRepository) CreateProfile(profileData req.ProfileReq) error {
	// Implement logic to create a new profile in MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := r.client.Collection(entity.Profiles{}.CollectionNameProfiles())
	convertdata, err := mapper.ToEntityProfile(&profileData)
	_, err = collection.InsertOne(ctx, convertdata)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProfileRepository) GetProfileByID(profileID string) (*entity.Profiles, error) {
	// Implement logic to retrieve a profile by ID from MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := r.client.Collection(entity.Profiles{}.CollectionNameProfiles())
	filter := bson.M{"_id": profileID}
	var profile entity.Profiles
	err := collection.FindOne(ctx, filter).Decode(&profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *ProfileRepository) GetProfileByUserID(userID string) (*entity.Profiles, error) {
	// Implement logic to retrieve a profile by UserID from MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := r.client.Collection(entity.Profiles{}.CollectionNameProfiles())
	filter := bson.M{"user_id": userID}
	var profile entity.Profiles
	err := collection.FindOne(ctx, filter).Decode(&profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *ProfileRepository) UpdateProfile(profileID string, updateData req.ProfileReq) error {
	dataprofile, err := r.GetProfileByID(profileID)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := r.client.Collection(entity.Profiles{}.CollectionNameProfiles())
	updatedEntity, err := mapper.ToEntityUpdateProfile(dataprofile, &updateData)
	filter := bson.M{"_id": profileID}
	_, err = collection.ReplaceOne(ctx, filter, updatedEntity)
	if err != nil {
		return err
	}
	return nil
}
