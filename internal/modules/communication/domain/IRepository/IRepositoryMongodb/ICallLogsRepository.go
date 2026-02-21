package IRepositoryMongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/entity"
)

type ICallLogsRepository interface {
	CreateCallLog(ctx context.Context, callLog *entity.CallLog) error
	CreateBulkCallLogs(ctx context.Context, callLogs []*entity.CallLog) (int64, []*mongodbErrors.BulkError, error)
	GetCallLogByID(ctx context.Context, id string) (*entity.CallLog, error)
	GetCallLogsByConversationID(ctx context.Context, conversationID string, cursor string, limit int) (*dto.PaginationRes, error)
	UpdateCallLog(ctx context.Context, callLog *entity.CallLog) error
	UpdateBulkCallLogs(ctx context.Context, callLogs []*entity.CallLog) (int64, []*mongodbErrors.BulkError, error)
	DeleteCallLog(ctx context.Context, id string) error
	DeleteBulkCallLogs(ctx context.Context, ids []string) (int64, []*mongodbErrors.BulkError, error)
}
