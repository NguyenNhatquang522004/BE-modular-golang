package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"
)

type CreateReelUseCase struct {
	// Define dependencies for creating a reel here (e.g., repositories, services)
	reelRepo IRepositoryMongodb.IReelRepository
}

func NewCreateReelUseCase( /* Add dependencies here */ ) *CreateReelUseCase {
	return &CreateReelUseCase{
		// Initialize dependencies here
	}
}
func (uc *CreateReelUseCase) Execute(ctx context.Context, req *req.CreateReelRequest) (*res.FailedReelResponse, error) {
	// Implement the logic for creating a reel here
	// Validate the request, interact with repositories/services, and return the appropriate response
	entity, err := mapper.ToEntityReel(req.ReelReq)
	if err != nil {
		return &res.FailedReelResponse{
			ReelID:       entity.ID.Hex(), // Assuming entity has an ID field of type primitive.ObjectID
			UserID:       req.UserID,
			ErrorMessage: err.Error(),
		}, err
	}

	return &res.FailedReelResponse{
		ReelID:       entity.ID.Hex(), // Assuming entity has an ID field of type primitive.ObjectID
		UserID:       req.UserID,
		ErrorMessage: "",
	}, nil
}
