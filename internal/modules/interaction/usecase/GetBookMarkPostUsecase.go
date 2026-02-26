package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/IRepository/IRepositoryMongoDB"
)

type GetBookMarkPostUsecase struct {
	saveItemRepo IRepositoryMongoDB.ISavedItemsRepository
}

func NewGetBookMarkPostUsecase(saveItemRepo IRepositoryMongoDB.ISavedItemsRepository) *GetBookMarkPostUsecase {
	return &GetBookMarkPostUsecase{
		saveItemRepo: saveItemRepo,
	}
}

func (u *GetBookMarkPostUsecase) Execute(ctx context.Context, req *req.GetBookMarkPostRequest) (*response.Response, error) {
	// Thực hiện logic để lấy danh sách bài viết đã bookmark dựa trên req.PaginationReq
	// Trả về response chứa danh sách bài viết đã bookmark
	data, err := u.saveItemRepo.PaginationSaveItem(ctx, req.UserID, req.Cursor, req.Limit)
	if err != nil {
		return response.NewResponse(
			response.WithData(nil),
			response.WithMessage("Lấy danh sách bài viết đã bookmark thất bại"),
			response.WithStatus("error"),
		), err
	}
	return response.NewResponse(
		response.WithData(data),
		response.WithMessage("Lấy danh sách bài viết đã bookmark thành công"),
		response.WithStatus("success"),
	), nil
}
