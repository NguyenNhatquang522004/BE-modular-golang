package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/res"
)

type ICreateGroup interface {
	Execute(ctx context.Context, req *req.CreateGroupRequest) (*res.FailedGroup, error)
}
type IUpdateGroup interface {
	Execute(ctx context.Context, req *req.UpdateGroupRequest) (*res.FailedGroup, error)
}
type IDeleteGroup interface {
	Execute(ctx context.Context, req *req.DeleteGroupRequest) ([]*res.FailedGroup, error)
}
type IAcceptGroupJoinRequest interface {
	Execute(ctx context.Context, req *req.AcceptGroupJoinRequest) (*res.FailedMember, error)
}
type IAddMemberGroup interface {
	Execute(ctx context.Context, req *req.AddMemberGroupRequest) (*res.FailedMember, error)
}
type IRemoveMemberGroup interface {
	Execute(ctx context.Context, req *req.RemoveMemberGroupRequest) (*res.FailedMember, error)
}
type IUpdateMemberGroup interface {
	Execute(ctx context.Context, req *req.UpdateMemberGroupRequest) (*res.FailedMember, error)
}
type ICreateGroupQA interface {
	Execute(ctx context.Context, req *req.CreateGroupQARequest) (*res.FailedGroupQA, error)
}
type IUpdateGroupQA interface {
	Execute(ctx context.Context, req *req.UpdateGroupQARequest) (*res.FailedGroupQA, error)
}
type IDeleteGroupQA interface {
	Execute(ctx context.Context, req *req.DeleteGroupQARequest) (*res.FailedGroupQA, error)
}
type ICreateGroupEvent interface {
	Execute(ctx context.Context, req *req.CreateGroupEventRequest) (*res.FailedGroupEvent, error)
}
type IUpdateGroupEvent interface {
	Execute(ctx context.Context, req *req.UpdateGroupEventRequest) (*res.FailedGroupEvent, error)
}
type IDeleteGroupEvent interface {
	Execute(ctx context.Context, req *req.DeleteGroupEventRequest) (*res.FailedGroupEvent, error)
}
type IStatsGroupEvent interface {
	Execute(ctx context.Context, req *req.StatsGroupEventRequest) (*res.FailedGroupEvent, error)
}
type ICreateGroupFile interface {
	Execute(ctx context.Context, req *req.CreateGroupFileRequest) (*res.FailGroupFile, error)
}
type IDeleteGroupFile interface {
	Execute(ctx context.Context, req *req.DeleteGroupFileRequest) (*res.FailGroupFile, error)
}
type IDownloadGroupFile interface {
	Execute(ctx context.Context, req *req.DownloadGroupFileRequest) (*res.FailGroupFile, error)
}

type Usecase struct{}

func NewUsecase() *Usecase {
	return &Usecase{}
}
