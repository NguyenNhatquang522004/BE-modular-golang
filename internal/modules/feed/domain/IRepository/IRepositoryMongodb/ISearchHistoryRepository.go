package IRepositoryMongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/feed/domain/entity"
)

type ISearchHistoryRepository interface {
	CreateSearchHistory(ctx context.Context, searchHistory *entity.SearchHistory) error
	CreateBulkSearchHistory(ctx context.Context, searchHistories []entity.SearchHistory) ([]*entity.SearchHistory, []*mongodbErrors.BulkError, error)
	GetSearchHistoriesByUserID(ctx context.Context, userID string, cursor string, limit int) (*dto.PaginationRes, error)
	GetSearchHistoryByUserIDTop(ctx context.Context, userID string) ([]*entity.SearchHistory, error)
	UpdateSearchHistory(ctx context.Context, searchHistory *entity.SearchHistory) error
	UpdateBulkSearchHistory(ctx context.Context, searchHistories []entity.SearchHistory) ([]*entity.SearchHistory, []*mongodbErrors.BulkError, error)
	DeleteSearchHistory(ctx context.Context, id string) error
	DeleteBulkSearchHistory(ctx context.Context, ids []string) ([]*entity.SearchHistory, []*mongodbErrors.BulkError, error)
}
