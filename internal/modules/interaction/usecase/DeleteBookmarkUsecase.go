package usecase

import (
	"context"
	"net/http"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/IRepository/IRepositoryMongoDB"
)

type DeleteBookmarkUsecase struct {
	saveItemRepo IRepositoryMongoDB.ISavedItemsRepository
}

func NewDeleteBookmarkUsecase(saveItemRepo IRepositoryMongoDB.ISavedItemsRepository) *DeleteBookmarkUsecase {
	return &DeleteBookmarkUsecase{
		saveItemRepo: saveItemRepo,
	}
}

func (u *DeleteBookmarkUsecase) Execute(ctx context.Context, req *req.DeleteBookmarkRequest) (*response.Response, error) {
	// Implement the logic for deleting a bookmark here

	err := u.saveItemRepo.DeleteSaveItem(ctx, req.SaveItemID)
	if err != nil {
		return nil, err
	}
	
	return response.NewResponse(response.WithData(""), response.WithMessage("Bookmark deleted successfully"), response.WithStatus(http.StatusOK)), nil
}
