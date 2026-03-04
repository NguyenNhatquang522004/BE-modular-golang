package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryPostgres"
)

type UpdateAdaccountUsecase struct {
	accountRepo IRepositoryPostgres.IAdAccountsRepository
}

func NewUpdateAdaccountUsecase(accountRepo IRepositoryPostgres.IAdAccountsRepository) *UpdateAdaccountUsecase {
	return &UpdateAdaccountUsecase{
		accountRepo: accountRepo,
	}
}
func (u *UpdateAdaccountUsecase) Execute(ctx context.Context, req *req.UpdateAdAccountRequest) (*res.FailedAdAccountResponse, error) {
	dataaccount, err := u.accountRepo.GetAdAccountByIDDetail(ctx, req.AccountID.String())
	if err != nil {
		return &res.FailedAdAccountResponse{
			AccountID:    req.AccountID.String(),
			UserActionID: req.UserActionID, // cần thêm trường UserActionID vào UpdateAdAccountRequest nếu muốn tracking
			ErrorMessage: "Failed to get ad account by ID: " + err.Error(),
		}, err
	}
	if dataaccount == nil {
		return &res.FailedAdAccountResponse{
			AccountID:    req.AccountID.String(),
			UserActionID: req.UserActionID, // cần thêm trường UserActionID vào UpdateAdAccountRequest nếu muốn tracking
			ErrorMessage: "Ad account not found",
		}, nil
	}
	if dataaccount.OwnerUserID.String() != req.UserActionID {
		return &res.FailedAdAccountResponse{
			AccountID:    req.AccountID.String(),
			UserActionID: req.UserActionID, // cần thêm trường UserActionID vào UpdateAdAccountRequest nếu muốn tracking
			ErrorMessage: "User does not have permission to update this ad account",
		}, nil
	}
	mapper.UpdateToEntityAdAccount(*req.AdAccountReq, dataaccount)
	err = u.accountRepo.UpdateAdAccount(ctx, dataaccount)
	if err != nil {
		return &res.FailedAdAccountResponse{
			AccountID:    req.AccountID.String(),
			UserActionID: req.UserActionID, // cần thêm trường UserActionID vào UpdateAdAccountRequest nếu muốn tracking
			ErrorMessage: "Failed to update ad account: " + err.Error(),
		}, err
	}
	return &res.FailedAdAccountResponse{
		AccountID:    req.AccountID.String(),
		UserActionID: req.UserActionID, // cần thêm trường UserActionID vào UpdateAdAccountRequest nếu muốn tracking
		ErrorMessage: "UpdateAdaccountUsecase is not implemented yet",
	}, nil
}
