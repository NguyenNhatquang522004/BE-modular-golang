package strategy

import (
	"context"
	"reflect"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryMongodb"
)

type EditLogStrategy struct {
	editLogRepo IRepositoryMongodb.IPostEditLogsRepository
}

func NewEditLogStrategy(editLogRepo IRepositoryMongodb.IPostEditLogsRepository) *EditLogStrategy {
	return &EditLogStrategy{
		editLogRepo: editLogRepo,
	}
}
func (s *EditLogStrategy) HandlePublishDelete(ctx context.Context, request *req.DeletePostRequest) error {
	// Implement the logic to handle post deletion
	err := s.editLogRepo.DeleteByTargetID(ctx, request.PostID.PostID)
	if err != nil {
		return err
	}
	return nil
}

func (s *EditLogStrategy) GetDeleteType() reflect.Type {
	return reflect.TypeOf(s.editLogRepo)
}
