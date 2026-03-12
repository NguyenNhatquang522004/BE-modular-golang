package mongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
	"go.mongodb.org/mongo-driver/mongo"
)

type ArtistRepository struct {
	client    *mongo.Database
	redisRepo IRepositoryShare.IRedis
}

func NewArtistRepository(client *mongo.Database, redisRepo IRepositoryShare.IRedis) *ArtistRepository {
	return &ArtistRepository{
		client:    client,
		redisRepo: redisRepo,
	}
}

func (r *ArtistRepository) CreatedArtist(ctx context.Context, artist *entity.Artist) error {
	collection := r.client.Collection(entity.Artist{}.CollectionName())
	_, err := collection.InsertOne(ctx, artist)
	return err

}
func (r *ArtistRepository) UpdatedArtist(ctx context.Context, artist *entity.Artist) error {
	collection := r.client.Collection(entity.Artist{}.CollectionName())
	_, err := collection.UpdateOne(ctx, map[string]interface{}{"_id": artist.ID}, map[string]interface{}{"$set": artist})
	return err
}
func (r *ArtistRepository) DeletedArtist(ctx context.Context, artistID string) error {
	collection := r.client.Collection(entity.Artist{}.CollectionName())
	_, err := collection.DeleteOne(ctx, map[string]interface{}{"_id": artistID})
	return err
}
func (r *ArtistRepository) GetArtistByID(ctx context.Context, artistID string) (*entity.Artist, error) {
	collection := r.client.Collection(entity.Artist{}.CollectionName())
	var artist entity.Artist
	err := collection.FindOne(ctx, map[string]interface{}{"_id": artistID}).Decode(&artist)
	if err != nil {
		return nil, err
	}
	return &artist, nil
}
