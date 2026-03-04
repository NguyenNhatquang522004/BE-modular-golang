package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryPostgres"
)

type UpdateBalanceUsecase struct {
	accountRepo IRepositoryPostgres.IAdAccountsRepository
}

func NewUpdateBalanceUsecase(accountRepo IRepositoryPostgres.IAdAccountsRepository) *UpdateBalanceUsecase {
	return &UpdateBalanceUsecase{
		accountRepo: accountRepo,
	}
}
func (u *UpdateBalanceUsecase) Execute(ctx context.Context, req *req.UpdateBalanceRequest) (*res.FailedAdAccountResponse, error) {
	// Logic để cập nhật số dư tài khoản quảng cáo
	// 1. Validate request
	// 2. Lấy thông tin tài khoản quảng cáo từ repository
	// 3. Cập nhật số dư dựa trên req.Amount
	// 4. Lưu thay đổi vào repository
	// 5. Trả về response hoặc lỗi nếu có
	dataaccount, err := u.accountRepo.GetAdAccountByIDDetail(ctx, req.AccountID)
	if err != nil {
		return &res.FailedAdAccountResponse{
			AccountID:    req.AccountID,
			UserActionID: req.UserActionID,
			ErrorMessage: "Failed to get ad account by ID: " + err.Error(),
		}, err
	}
	if dataaccount == nil {
		return &res.FailedAdAccountResponse{
			AccountID:    req.AccountID,
			UserActionID: req.UserActionID,
			ErrorMessage: "Ad account not found",
		}, nil
	}
	if dataaccount.OwnerUserID.String() != req.UserActionID {
		return &res.FailedAdAccountResponse{
			AccountID:    req.AccountID,
			UserActionID: req.UserActionID,
			ErrorMessage: "User does not have permission to update this ad account",
		}, nil
	}
	// Cập nhật số dư
	dataaccount.Balance += req.Amount
	err = u.accountRepo.UpdateAdAccount(ctx, dataaccount)
	if err != nil {
		return &res.FailedAdAccountResponse{
			AccountID:    req.AccountID,
			UserActionID: req.UserActionID,
			ErrorMessage: "Failed to update ad account balance: " + err.Error(),
		}, err
	}
	return &res.FailedAdAccountResponse{
		AccountID:    req.AccountID,
		UserActionID: req.UserActionID,
		ErrorMessage: "Ad account balance updated successfully",
	}, nil
}
