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
type IAddMemberGroup interface {
	Execute(ctx context.Context, req *req.AddMemberGroupRequest) (*res.FailedMember, error)
}
type IRemoveMemberGroup interface {
	Execute(ctx context.Context) error
}
type ICreateGroupQA interface {
	Execute(ctx context.Context) error
}
type IUpdateGroupQA interface {
	Execute(ctx context.Context) error
}
type IDeleteGroupQA interface {
	Execute(ctx context.Context) error
}
type ICreateGroupEvent interface {
	Execute(ctx context.Context) error
}
type IUpdateGroupEvent interface {
	Execute(ctx context.Context) error
}
type IDeleteGroupEvent interface {
	Execute(ctx context.Context) error
}
type ICreateGroupFile interface {
	Execute(ctx context.Context) error
}
type IDeleteGroupFile interface {
	Execute(ctx context.Context) error
}
type IDownloadGroupFile interface {
	Execute(ctx context.Context) error
}

type Usecase struct{}

func NewUsecase() *Usecase {
	return &Usecase{}
}
