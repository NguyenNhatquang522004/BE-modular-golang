package usecase

import "go.mongodb.org/mongo-driver/mongo"

type ProfileUseCase struct {
	client *mongo.Database
}

func NewProfileUseCase(client *mongo.Database) *ProfileUseCase {
	return &ProfileUseCase{
		client: client,
	}
}


