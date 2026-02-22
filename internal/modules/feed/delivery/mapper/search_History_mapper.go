package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/feed/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/feed/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/feed/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- REQUEST TO ENTITY ---

func ToEntitySearchHistory(dto req.SearchHistoryReq) entity.SearchHistory {
	return entity.SearchHistory{
		ID:         primitive.NewObjectID(), // Tạo ID mới cho bản ghi mới
		UserID:     dto.UserID,
		Keyword:    dto.Keyword,
		TargetID:   dto.TargetID,
		TargetType: dto.TargetType,
		CreatedAt:  time.Now(),
	}
}

func UpdateToEntitySearchHistory(dto req.SearchHistoryReq, existingEntity *entity.SearchHistory) {
	existingEntity.UserID = dto.UserID
	existingEntity.Keyword = dto.Keyword
	existingEntity.TargetID = dto.TargetID
	existingEntity.TargetType = dto.TargetType
	// CreatedAt thường không đổi khi update, trừ khi bạn có field UpdatedAt
}

// --- ENTITY TO RESPONSE ---

func ToResponseSearchHistory(ent entity.SearchHistory) res.SearchHistoryRes {
	return res.SearchHistoryRes{
		ID:         ent.ID.Hex(), // Chuyển ObjectID sang String
		UserID:     ent.UserID,
		Keyword:    ent.Keyword,
		TargetID:   ent.TargetID,
		TargetType: ent.TargetType,
		CreatedAt:  ent.CreatedAt,
	}
}

// ToResponseListSearchHistory hỗ trợ chuyển đổi một mảng kết quả từ MongoDB
func ToResponseListSearchHistory(entities []entity.SearchHistory) []res.SearchHistoryRes {
	responses := make([]res.SearchHistoryRes, len(entities))
	for i, ent := range entities {
		responses[i] = ToResponseSearchHistory(ent)
	}
	return responses
}
