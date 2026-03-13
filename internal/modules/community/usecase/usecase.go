package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/res"
)

type ICreateGroupUsecase interface {
	Execute(ctx context.Context, req *req.CreateGroupRequest) (*res.FailedGroup, error)
}
type IUpdateGroupUsecase interface {
	Execute(ctx context.Context, req *req.UpdateGroupRequest) (*res.FailedGroup, error)
}
type IDeleteGroupUsecase interface {
	Execute(ctx context.Context, req *req.DeleteGroupRequest) ([]*res.FailedGroup, error)
}
type IAcceptGroupJoinUsecase interface {
	Execute(ctx context.Context, req *req.AcceptGroupJoinRequest) (*res.FailedMember, error)
}
type IAddMemberGroupUsecase interface {
	Execute(ctx context.Context, req *req.AddMemberGroupRequest) (*res.FailedMember, error)
}
type IRemoveMemberGroupUsecase interface {
	Execute(ctx context.Context, req *req.RemoveMemberGroupRequest) (*res.FailedMember, error)
}
type IUpdateMemberGroupUsecase interface {
	Execute(ctx context.Context, req *req.UpdateMemberGroupRequest) (*res.FailedMember, error)
}
type ICreateGroupQAUsecase interface {
	Execute(ctx context.Context, req *req.CreateGroupQARequest) (*res.FailedGroupQA, error)
}
type IUpdateGroupQAUsecase interface {
	Execute(ctx context.Context, req *req.UpdateGroupQARequest) (*res.FailedGroupQA, error)
}
type IDeleteGroupQAUsecase interface {
	Execute(ctx context.Context, req *req.DeleteGroupQARequest) (*res.FailedGroupQA, error)
}
type ICreateGroupEventUsecase interface {
	Execute(ctx context.Context, req *req.CreateGroupEventRequest) (*res.FailedGroupEvent, error)
}
type IUpdateGroupEventUsecase interface {
	Execute(ctx context.Context, req *req.UpdateGroupEventRequest) (*res.FailedGroupEvent, error)
}
type IDeleteGroupEventUsecase interface {
	Execute(ctx context.Context, req *req.DeleteGroupEventRequest) (*res.FailedGroupEvent, error)
}
type IStatsGroupEventUsecase interface {
	Execute(ctx context.Context, req *req.StatsGroupEventRequest) (*res.FailedGroupEvent, error)
}
type ICreateGroupFileUsecase interface {
	Execute(ctx context.Context, req *req.CreateGroupFileRequest) (*res.FailGroupFile, error)
}
type IDeleteGroupFileUsecase interface {
	Execute(ctx context.Context, req *req.DeleteGroupFileRequest) (*res.FailGroupFile, error)
}

type Usecase struct{}

func NewUsecase() *Usecase {
	return &Usecase{}
}
