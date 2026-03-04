package IRepositoryPostgres

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/postgresErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/entity"
)

type IAdAccountsRepository interface {
	// Define the methods for the AdAccountsRepository interface here
	CreateAdAccount(ctx context.Context, account *entity.AdAccount) error
	CreateBulkAdAccounts(ctx context.Context, accounts []*entity.AdAccount) (int64, []*postgresErrors.AdAccountsBulkError, error)
	GetAdAccountByIDDetail(ctx context.Context, accountID string) (*entity.AdAccount, error)
	GetAdAccountByID(ctx context.Context, cursor string, limit int) (*dto.PaginationRes, error)
	GetAdAccountsByOwnerUserIDDetail(ctx context.Context, ownerUserID string) (*entity.AdAccount, error)
	GetAdAccountsByOwnerUserID(ctx context.Context, ownerUserID string, cursor string, limit int) (*dto.PaginationRes, error)
	UpdateAdAccount(ctx context.Context, account *entity.AdAccount) error
	UpdateBulkAdAccounts(ctx context.Context, accounts []*entity.AdAccount) (int64, []*postgresErrors.AdAccountsBulkError, error)
	DeleteAdAccount(ctx context.Context, accountID string) error
	DeleteBulkAdAccounts(ctx context.Context, accountIDs []string) (int64, []*postgresErrors.AdAccountsBulkError, error)
}
