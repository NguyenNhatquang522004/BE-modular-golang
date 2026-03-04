package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryPostgres"
)

type CreateAdaccountUsecase struct {
	accountRepo IRepositoryPostgres.IAdAccountsRepository
}

func NewCreateAdaccountUsecase(accountRepo IRepositoryPostgres.IAdAccountsRepository) *CreateAdaccountUsecase {
	return &CreateAdaccountUsecase{

		accountRepo: accountRepo,
	}
}

func (u *CreateAdaccountUsecase) Execute(ctx context.Context, req *req.CreateAdAccountRequest) (*res.FailedAdAccountResponse, error) {
	dataaccount, err := u.accountRepo.GetAdAccountsByOwnerUserIDDetail(ctx, req.OwnerUserID.String())
	if err != nil {
		return &res.FailedAdAccountResponse{
			AccountID:    req.AccountID.String(),
			UserActionID: "", // cần thêm trường UserActionID vào CreateAdAccountRequest nếu muốn tracking
			ErrorMessage: "Failed to get ad accounts for user: " + err.Error(),
		}, err
	}
	if dataaccount != nil {
		return &res.FailedAdAccountResponse{
			AccountID:    req.AccountID.String(),
			UserActionID: "", // cần thêm trường UserActionID vào CreateAdAccountRequest nếu muốn tracking
			ErrorMessage: "User already has ad accounts, cannot create another one",
		}, nil
	}
	entity := mapper.ToEntityAdAccount(*req.AdAccountReq)
	if err := u.accountRepo.CreateAdAccount(ctx, &entity); err != nil {
		return &res.FailedAdAccountResponse{
			AccountID:    req.AccountID.String(),
			UserActionID: "", // cần thêm trường UserActionID vào CreateAdAccountRequest nếu muốn tracking
			ErrorMessage: "Failed to create ad account: " + err.Error(),
		}, err
	}
	return &res.FailedAdAccountResponse{
		AccountID:    req.AccountID.String(),
		UserActionID: "", // cần thêm trường UserActionID vào CreateAdAccountRequest nếu muốn tracking
		ErrorMessage: "Ad account created successfully",
	}, nil
}
