package IRepositoryMongoDB

import (
	"context"

	mongodb "github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/entity"
)

type ISavedItemsRepository interface {
	CreateSaveItem(ctx context.Context, saveItem *entity.UserSavedItem) error
	CreateBulkSaveItem(ctx context.Context, saveItems []*entity.UserSavedItem) (int64, []*mongodb.BulkError, error)
	GetSaveItemByID(ctx context.Context, saveItemID string) (*entity.UserSavedItem, error)
	GetBulkSaveItemByID(ctx context.Context, saveItemIDs []string) ([]*entity.UserSavedItem, error)
	GetSaveItemsByUserID(ctx context.Context, userID string) ([]*entity.UserSavedItem, error)
	UpdateSaveItem(ctx context.Context, saveItem *entity.UserSavedItem) error
	UpdateBulkSaveItem(ctx context.Context, saveItems []*entity.UserSavedItem) (int64, []*mongodb.BulkError, error)
	DeleteSaveItem(ctx context.Context, saveItemID string) error
	DeleteBulkSaveItem(ctx context.Context, saveItemIDs []string) (int64, []*mongodb.BulkError, error)
	PaginationSaveItem(ctx context.Context, userID string, cursor string, limit int) (*dto.PaginationRes, error)
}
