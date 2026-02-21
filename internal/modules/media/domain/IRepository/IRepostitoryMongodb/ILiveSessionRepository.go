package IRepostitoryMongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
)

type ILiveSessionRepository interface {
	// Define methods for LiveSessionRepository here
	CreateLiveSession(ctx context.Context, liveSession *entity.LiveSession) error
	CreateBulkLiveSessions(ctx context.Context, liveSessions []*entity.LiveSession) (int64, []*mongodbErrors.BulkError, error)
	GetLiveSessionByID(ctx context.Context, id string) (*entity.LiveSession, error)
	GetLiveSessionsByUserID(ctx context.Context, userID string, cursor string, limit int) (*dto.PaginationRes, error)
	GetLiveSessionsByCategoryID(ctx context.Context, categoryID string, cursor string, limit int) (*dto.PaginationRes, error)
	UpdateLiveSession(ctx context.Context, liveSession *entity.LiveSession) error
	UpdateBulkLiveSessions(ctx context.Context, liveSessions []*entity.LiveSession) (int64, []*mongodbErrors.BulkError, error)
	DeleteLiveSession(ctx context.Context, id string) error
	DeleteBulkLiveSessions(ctx context.Context, ids []string) (int64, []*mongodbErrors.BulkError, error)
	CheckLiveSessionExists(ctx context.Context, id string) (bool, error)
	CheckLiveSessionExistsByStreamKey(ctx context.Context, HostUserID string, streamKey string) (bool, error)
}
