package IRepositoryMongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/entity"
)

type IGroupjoinQuestionsRepository interface {
	CreateGroupJoinQuestion(ctx context.Context, question *entity.GroupJoinQuestion) error
	CreateBulkGroupJoinQuestions(ctx context.Context, questions []entity.GroupJoinQuestion) (int64, []*dto.BulkError, error)
	GetGroupJoinQuestionByID(ctx context.Context, questionID string) (*entity.GroupJoinQuestion, error)
	GetGroupJoinQuestionsByGroupID(ctx context.Context, groupID string, cursor string, limit int) (*dto.PaginationRes, error)
	UpdateGroupJoinQuestion(ctx context.Context, question *entity.GroupJoinQuestion) error
	UpdateBulkGroupJoinQuestions(ctx context.Context, questions []entity.GroupJoinQuestion) (int64, []*dto.BulkError, error)
	DeleteGroupJoinQuestion(ctx context.Context, questionID string) error
	DeleteBulkGroupJoinQuestions(ctx context.Context, questionIDs []string) (int64, []*dto.BulkError, error)
	DeleteGroupJoinQuestionsByGroupID(ctx context.Context, groupID string) (int64, error)
	DeleteBulkGroupJoinQuestionsByGroupID(ctx context.Context, groupIDs []string) (int64, []*dto.BulkError, error)
	DeleteGroupJoinQuestionByIDAndGroupID(ctx context.Context, questionID string, groupID string) error
	DeleteBulkGroupJoinQuestionByIDAndGroupID(ctx context.Context, questionIDs []string, groupID string) (int64, []*dto.BulkError, error)
}
