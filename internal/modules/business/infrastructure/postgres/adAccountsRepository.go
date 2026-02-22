package postgres

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/postgresErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdAccountsRepository struct {
	db        *gorm.DB
	redisRepo IRepositoryShare.IRedis
}

// Implement the methods for the AdAccountsRepository struct here
func NewAdAccountsRepository(db *gorm.DB, redisRepo IRepositoryShare.IRedis) *AdAccountsRepository {
	return &AdAccountsRepository{
		db:        db,
		redisRepo: redisRepo,
	}
}
func (r *AdAccountsRepository) CreateAdAccount(ctx context.Context, account *entity.AdAccount) error {
	return r.db.WithContext(ctx).Create(account).Error
}
func (r *AdAccountsRepository) CreateBulkAdAccounts(ctx context.Context, accounts []*entity.AdAccount) (int64, []*postgresErrors.AdAccountsBulkError, error) {
	var bulkErrors []*postgresErrors.AdAccountsBulkError
	for _, account := range accounts {
		if err := r.db.WithContext(ctx).Create(account).Error; err != nil {
			bulkErrors = append(bulkErrors, &postgresErrors.AdAccountsBulkError{
				ID:          account.ID.String(),
				OwnerUserID: account.OwnerUserID.String(),
				Reason:      err.Error(),
			})
		}
	}
	return int64(len(accounts) - len(bulkErrors)), bulkErrors, nil
}
func (r *AdAccountsRepository) GetAdAccountByIDDetail(ctx context.Context, accountID string) (*entity.AdAccount, error) {
	var account *entity.AdAccount
	if err := r.db.WithContext(ctx).Where(&entity.AdAccount{ID: uuid.MustParse(accountID)}).First(&account).Error; err != nil {
		return nil, err
	}
	return account, nil
}
func (r *AdAccountsRepository) GetAdAccountByID(ctx context.Context, cursor string, limit int) (*dto.PaginationRes, error) {
	if cursor == "" {
		itemsSet := []string{
			"adAccount_cache_bulk",
			"adAccount_cache_bulk_nextcursor",
			"adAccount_cache_bulk_hasnext",
			"adAccount_cache_bulk_limit",
		}
		data, nextcursor, hasnext, limit, err := r.redisRepo.CustomizeGetCache(context.Background(), itemsSet)
		if err == nil {
			return &dto.PaginationRes{
				Limit:      limit,
				NextCursor: nextcursor,
				HasNext:    hasnext,
				Data:       data,
			}, nil
		}
		return &dto.PaginationRes{
			NextCursor: nextcursor,
			HasNext:    hasnext,
			Data:       data,
			Limit:      limit,
		}, nil
	}
	var accounts []*entity.AdAccount
	querylimit := limit + 1
	query := r.db.WithContext(ctx).Model(&entity.AdAccount{}).Order("created_at DESC , id DESC").Limit(querylimit)
	if cursor != "" {
		lastCreatedAt, lastID, err := utils.DecodeCursor(cursor)
		if err != nil {
			return nil, err
		}
		query = query.Where("(created_at < ?) OR (created_at = ? AND id < ?)", lastCreatedAt, lastCreatedAt, lastID)
	}
	err := query.Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	var nextCursor string
	if len(accounts) == querylimit {
		lastAccount := accounts[limit-1]
		nextCursor = utils.EncodeCursor(lastAccount.CreatedAt, lastAccount.ID)
		accounts = accounts[:limit]
	}
	var hasNext bool
	if len(accounts) == querylimit {
		hasNext = true
	}
	if cursor == "" {
		itemsSet := map[string]any{
			"adAccount_cache_bulk":            accounts,
			"adAccount_cache_bulk_nextcursor": nextCursor,
			"adAccount_cache_bulk_hasnext":    hasNext,
			"adAccount_cache_bulk_limit":      limit,
		}
		err = r.redisRepo.CustomizeSetCache(context.Background(), itemsSet)
		if err != nil {
			return nil, err
		}
	}
	return &dto.PaginationRes{
		Data:       accounts,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *AdAccountsRepository) GetAdAccountsByOwnerUserID(ctx context.Context, ownerUserID string, cursor string, limit int) (*dto.PaginationRes, error) {
	if cursor == "" {
		itemsSet := []string{
			"adAccount_cache_bulk_ownerUserID_" + ownerUserID,
			"adAccount_cache_bulk_nextcursor_ownerUserID_" + ownerUserID,
			"adAccount_cache_bulk_hasnext_ownerUserID_" + ownerUserID,
			"adAccount_cache_bulk_limit_ownerUserID_" + ownerUserID,
		}
		data, nextcursor, hasnext, limit, err := r.redisRepo.CustomizeGetCache(context.Background(), itemsSet)
		if err == nil {
			return &dto.PaginationRes{
				Limit:      limit,
				NextCursor: nextcursor,
				HasNext:    hasnext,
				Data:       data,
			}, nil
		}
		return &dto.PaginationRes{
			NextCursor: nextcursor,
			HasNext:    hasnext,
			Data:       data,
			Limit:      limit,
		}, nil
	}
	var accounts []*entity.AdAccount
	querylimit := limit + 1
	query := r.db.WithContext(ctx).Model(&entity.AdAccount{}).Where(&entity.AdAccount{OwnerUserID: uuid.MustParse(ownerUserID)}).Order("created_at DESC , id DESC").Limit(querylimit)
	if cursor != "" {
		lastCreatedAt, lastID, err := utils.DecodeCursor(cursor)
		if err != nil {
			return nil, err
		}
		query = query.Where("(created_at < ?) OR (created_at = ? AND id < ?)", lastCreatedAt, lastCreatedAt, lastID)
	}
	err := query.Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	var nextCursor string
	if len(accounts) == querylimit {
		lastAccount := accounts[limit-1]
		nextCursor = utils.EncodeCursor(lastAccount.CreatedAt, lastAccount.ID)
		accounts = accounts[:limit]
	}
	var hasNext bool
	if len(accounts) == querylimit {
		hasNext = true
	}
	if cursor == "" {
		itemsSet := map[string]any{
			"adAccount_cache_bulk_ownerUserID_" + ownerUserID:            accounts,
			"adAccount_cache_bulk_nextcursor_ownerUserID_" + ownerUserID: nextCursor,
			"adAccount_cache_bulk_hasnext_ownerUserID_" + ownerUserID:    hasNext,
			"adAccount_cache_bulk_limit_ownerUserID_" + ownerUserID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(context.Background(), itemsSet)
		if err != nil {
			return nil, err
		}
	}
	return &dto.PaginationRes{
		Data:       accounts,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *AdAccountsRepository) UpdateAdAccount(ctx context.Context, account *entity.AdAccount) error {
	return r.db.WithContext(ctx).Model(&entity.AdAccount{}).Where(&entity.AdAccount{ID: account.ID}).Updates(account).Error
}
func (r *AdAccountsRepository) UpdateBulkAdAccounts(ctx context.Context, accounts []*entity.AdAccount) (int64, []*postgresErrors.AdAccountsBulkError, error) {
	var bulkErrors []*postgresErrors.AdAccountsBulkError
	for _, account := range accounts {
		if err := r.db.WithContext(ctx).Model(&entity.AdAccount{}).Where(&entity.AdAccount{ID: account.ID}).Updates(account).Error; err != nil {
			bulkErrors = append(bulkErrors, &postgresErrors.AdAccountsBulkError{
				ID:          account.ID.String(),
				OwnerUserID: account.OwnerUserID.String(),
				Reason:      err.Error(),
			})
		}
	}
	return int64(len(accounts) - len(bulkErrors)), bulkErrors, nil
}
func (r *AdAccountsRepository) DeleteAdAccount(ctx context.Context, accountID string) error {
	return r.db.WithContext(ctx).Where(&entity.AdAccount{ID: uuid.MustParse(accountID)}).Delete(&entity.AdAccount{}).Error
}
func (r *AdAccountsRepository) DeleteBulkAdAccounts(ctx context.Context, accountIDs []string) (int64, []*postgresErrors.AdAccountsBulkError, error) {
	var bulkErrors []*postgresErrors.AdAccountsBulkError
	for _, accountID := range accountIDs {
		if err := r.db.WithContext(ctx).Where(&entity.AdAccount{ID: uuid.MustParse(accountID)}).Delete(&entity.AdAccount{}).Error; err != nil {
			bulkErrors = append(bulkErrors, &postgresErrors.AdAccountsBulkError{
				ID:          accountID,
				OwnerUserID: "", // Không có thông tin OwnerUserID khi xóa, có thể để trống hoặc tìm kiếm trước khi xóa nếu cần
				Reason:      err.Error(),
			})
		}
	}
	return int64(len(accountIDs) - len(bulkErrors)), bulkErrors, nil
}
