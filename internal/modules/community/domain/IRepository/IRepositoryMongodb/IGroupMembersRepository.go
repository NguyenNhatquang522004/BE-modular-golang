package IRepositoryMongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/entity"
)

type IGroupMembersRepository interface {
	CreateGroupMember(ctx context.Context, groupMember *entity.GroupMember) error
	CreateBulkGroupMembers(ctx context.Context, groupMembers []*entity.GroupMember) (int64, []*dto.BulkError, error)
	GetGroupMemberByID(ctx context.Context, groupID string, cursor string, limit int) (*dto.PaginationRes, error)
	GetGroupMemberByUserID(ctx context.Context, userID string, cursor string, limit int) (*dto.PaginationRes, error)
	GetGroupMemberByUserIDAndGroupID(ctx context.Context, userID string, groupID string) (*entity.GroupMember, error)
	UpdateGroupMember(ctx context.Context, groupMember *entity.GroupMember) error
	UpdateBulkGroupMembers(ctx context.Context, groupMembers []*entity.GroupMember) (int64, []*dto.BulkError, error)
	DeleteGroupMember(ctx context.Context, groupID string) error
	DeleteBulkGroupMembers(ctx context.Context, groupIDs []string) (int64, []*dto.BulkError, error)
	DeleteGroupMemberByUserIDAndGroupID(ctx context.Context, userID string, groupID string) error
	DeleteBulkGroupMemberByUserIDAndGroupID(ctx context.Context, userID string, groupIDs []string) (int64, []*dto.BulkError, error)
	DeleteBulkGroupMemberByManyUserIDAndGroupID(ctx context.Context, userID []string, groupIDs []string) (int64, []*dto.BulkError, error)
}
