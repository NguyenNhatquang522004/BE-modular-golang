package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryMongodb"
)

type EditPostUseCase struct {
	editlog IRepositoryMongodb.IPostEditLogsRepository
}

func NewEditPostUseCase(editlog IRepositoryMongodb.IPostEditLogsRepository) *EditPostUseCase {
	return &EditPostUseCase{
		editlog: editlog,
	}
}

func (uc *EditPostUseCase) Execute(ctx context.Context, req *req.EditPostRequest) (*response.Response, error) {
	// 1. Tạo một bản ghi edit log mới với thông tin về phiên bản mới của bài viết, người chỉnh sửa, thời gian chỉnh sửa, v.v.
	data, err := uc.editlog.GetLatestVersion(ctx, req.EditLog.TargetID)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage(""),
			response.WithStatus("")), err
	}
	newVersion := data + 1
	req.EditLog.Version = newVersion
	data1 := mapper.ToPostEntityEditLogEntity(req.EditLog)
	err = uc.editlog.CreatePostEditLog(ctx, data1)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage(""),
			response.WithStatus("")), err
	}
	// 2. Trả về phản hồi thành công hoặc lỗi nếu có vấn đề xảy ra trong quá trình tạo bản ghi edit log.
	return response.NewResponse(response.WithData(""),
		response.WithMessage("Edit log created successfully"),
		response.WithStatus("success")), nil
}
