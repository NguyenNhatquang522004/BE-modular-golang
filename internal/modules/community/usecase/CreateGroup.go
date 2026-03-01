package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/IRepository/IRepositoryMongodb"
)

type CreateGroup struct {
	groupRepo IRepositoryMongodb.IGroupRepository
}

func NewCreateGroup(groupRepo IRepositoryMongodb.IGroupRepository) *CreateGroup {
	return &CreateGroup{
		groupRepo: groupRepo,
	}
}

func (c *CreateGroup) Execute(ctx context.Context, req *req.CreateGroupRequest) (*res.FailedGroup, error) {
	// Implement the logic to create a group here
	// This may involve validating the request, interacting with the repository, and returning an appropriate response
	entity := mapper.ToEntityGroup(req.CreateGroupReq)
	err := c.groupRepo.CreateGroup(ctx, entity)
	if err != nil {
		return &res.FailedGroup{
			GroupID:      entity.ID.Hex(), // Assuming the group ID is generated and available after creation
			UserID:       req.CreatorID,
			ErrorMessage: err,
		}, nil
	}
	return &res.FailedGroup{
		GroupID:      entity.ID.Hex(), // Replace with the actual group ID
		UserID:       req.CreatorID,
		ErrorMessage: nil,
	}, nil
}
