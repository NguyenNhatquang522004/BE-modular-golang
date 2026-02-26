package usecase

import (
	"context"
	"net/http"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/IRepository/IRepositoryMongoDB"
)

type BookmarkPostUsecase struct {
	saveItemRepo IRepositoryMongoDB.ISavedItemsRepository
}

func NewBookmarkPostUsecase(saveItemRepo IRepositoryMongoDB.ISavedItemsRepository) *BookmarkPostUsecase {
	return &BookmarkPostUsecase{
		saveItemRepo: saveItemRepo,
	}
}

func (u *BookmarkPostUsecase) Execute(ctx context.Context, req *req.BookmarkPostRequest) (*response.Response, error) {
	// Implement the logic for bookmarking a post here
	entity, err := mapper.ToEntityUserSavedItem(req.UserSavedItemReq)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus(http.StatusBadRequest)), err
	}
	err = u.saveItemRepo.CreateSaveItem(ctx, entity)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus(http.StatusBadRequest)), err
	}
	return response.NewResponse(response.WithData(""), response.WithMessage("Bookmark created successfully"), response.WithStatus(http.StatusOK)), nil
}
