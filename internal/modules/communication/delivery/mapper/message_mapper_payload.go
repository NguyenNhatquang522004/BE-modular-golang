// File: mapper/message_mapper.go
package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/communicationEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/entity"
	"github.com/gocql/gocql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func generateBucket(t time.Time) int {
	// Trả về định dạng 202603 cho tháng 3 năm 2026
	return t.Year()*100 + int(t.Month())
}
func ToMessageEntity(req *communicationEvent.CreateMessagePayload) (*entity.Message, []string) {
	now := time.Now()

	// 1. Xử lý MessageID (Sử dụng Parse để đảm bảo tính hợp lệ của TimeUUID)
	messageID, err := gocql.ParseUUID(req.MessageID)
	if err != nil {
		messageID = gocql.TimeUUID()
	}

	// 2. Trích xuất danh sách Asset IDs để lưu vào Cassandra
	// Chúng ta lưu ID để giữ bảng Message nhẹ nhàng
	assetIDs := make([]string, 0, len(req.Attachments))
	for _, att := range req.Attachments {
		att.AssetID = primitive.NewObjectID().Hex() // Đảm bảo mỗi Attachment có một ID duy nhất
		assetIDs = append(assetIDs, att.AssetID)
	}

	return &entity.Message{
		ConversationID:   req.ConversationID,
		Bucket:           generateBucket(now),
		MessageID:        messageID,
		SenderID:         req.SenderID,
		Type:             req.Type,
		Content:          req.Content,
		Attachments:      assetIDs, // Chỉ lưu IDs vào Cassandra
		IsEdited:         false,
		ReplyToMessageID: req.ReplyToMessageID,
		StoryRefID:       req.StoryRefID,
		IsRevoked:        false,
		CreatedAt:        now,
	}, assetIDs
}

func UpdateMessageMapper(existing *entity.Message, req *communicationEvent.UpdateMessagePayload) (*entity.Message, []communicationEvent.AttachmentPayload, []string) {
	var addedAssets []communicationEvent.AttachmentPayload
	var removedIDs []string
	hasChanged := false

	// 1. Xử lý Content
	if req.Content != nil && *req.Content != existing.Content {
		existing.Content = *req.Content
		hasChanged = true
	}

	// 2. Xử lý Attachments
	if req.Attachments != nil {
		newAssetIDs := make([]string, 0, len(*req.Attachments))

		// Tạo map để kiểm tra xem ID nào đã tồn tại trong Entity cũ chưa
		oldAssetsMap := make(map[string]bool)
		for _, id := range existing.Attachments {
			oldAssetsMap[id] = true
		}

		// DUYỆT BẰNG INDEX (i) thay vì copy value (att) để có thể thay đổi giá trị gốc
		for i := range *req.Attachments {
			// Lấy con trỏ của phần tử hiện tại
			att := &(*req.Attachments)[i]

			// TÍNH NĂNG MỚI: NẾU CHƯA CÓ ID -> TỰ ĐỘNG TẠO ID MỚI
			if att.AssetID == "" {
				att.AssetID = primitive.NewObjectID().Hex() // Tạo ID chuẩn MongoDB
			}

			newAssetIDs = append(newAssetIDs, att.AssetID)

			// Nếu AssetID này không có trong Entity cũ -> Đây là hàng mới thêm
			if !oldAssetsMap[att.AssetID] {
				// Vì att là con trỏ, ta lấy giá trị thực tế của nó để đưa vào mảng addedAssets
				addedAssets = append(addedAssets, *att)
			}
		}

		existing.Attachments = newAssetIDs
		hasChanged = true
	}

	// 3. Xử lý RemoveAttachmentIDs (Ghi nhận để xóa ở MongoDB)
	if len(req.RemoveAttachmentIDs) > 0 {
		removedIDs = req.RemoveAttachmentIDs
		// Thực hiện xóa các ID này khỏi Entity hiện tại (nếu Client không gửi list Attachments mới)
		if req.Attachments == nil {
			remainingIDs := []string{}
			toRemove := make(map[string]bool)
			for _, id := range req.RemoveAttachmentIDs {
				toRemove[id] = true
			}
			for _, id := range existing.Attachments {
				if !toRemove[id] {
					remainingIDs = append(remainingIDs, id)
				}
			}
			existing.Attachments = remainingIDs
		}
		hasChanged = true
	}

	// 4. Xử lý Revoke (Trọng yếu)
	if req.IsRevoked != nil && *req.IsRevoked {
		existing.IsRevoked = true
		existing.Content = "Tin nhắn đã bị thu hồi"
		existing.Attachments = []string{}
		existing.IsEdited = false // Thu hồi thì không hiện "Đã chỉnh sửa"
		return existing, nil, removedIDs
	}

	// 5. Cập nhật flag IsEdited
	if hasChanged {
		existing.IsEdited = true
	}

	return existing, addedAssets, removedIDs
}
