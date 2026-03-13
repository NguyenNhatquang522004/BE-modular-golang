package mapper

import (
	"errors"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/communityEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func parseOptionalObjectID(hex *string) *primitive.ObjectID {
	if hex == nil {
		return nil
	}
	id, err := primitive.ObjectIDFromHex(*hex)
	if err != nil {
		return nil
	}
	return &id
}

// MapCreatePayloadToEntity chuyển đổi Payload tạo file sang Entity chuẩn.
// Trả về error nếu GroupID không phải là chuẩn ObjectID của MongoDB.
func MapCreatePayloadToEntity(req *communityEvent.CreateGroupFilePayload) (*entity.GroupFile, error) {
	// Bắt buộc parse GroupID
	groupID, err := primitive.ObjectIDFromHex(req.GroupID)
	if err != nil {
		return nil, err
	}
	if req.ID == "" {
		return nil, errors.New("ID is required for tracking purposes, even if it will be ignored by DB")
	}
	ID, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return nil, err
	}
	now := time.Now()

	fileEntity := &entity.GroupFile{
		ID:            ID, // Tự động sinh ID mới
		GroupID:       groupID,
		UploaderID:    req.UploaderID, // Truyền từ ngoài vào (Middleware/JWT)
		FileName:      req.FileName,
		FileType:      req.FileType,
		FileSize:      req.FileSize,
		StorageFileID: req.StorageFileID,
		OriginalURL:   req.OriginalURL,
		ThumbnailURL:  req.ThumbnailURL,
		Caption:       req.Caption,
		DownloadCount: 0,   // Luôn khởi tạo bằng 0
		CreatedAt:     now, // Gán thời gian hiện tại
		UpdatedAt:     now,
		Metadata: entity.FileMetadata{
			Extension: req.Metadata.Extension,
			MimeType:  req.Metadata.MimeType,
			SizeBytes: req.Metadata.SizeBytes,
			Width:     req.Metadata.Width,
			Height:    req.Metadata.Height,
			Duration:  req.Metadata.Duration,
		},
	}

	// Xử lý các Optional ID (Nếu Client có gửi lên thì mới parse và gán)
	fileEntity.PostID = parseOptionalObjectID(req.PostID)
	fileEntity.CommentID = parseOptionalObjectID(req.CommentID)
	fileEntity.MessageID = parseOptionalObjectID(req.MessageID)

	return fileEntity, nil
}

// MapUpdatePayloadToEntity áp dụng các thay đổi từ payload vào entity có sẵn.
// Trả về chính entity đó sau khi đã được cập nhật.
func MapUpdatePayloadToEntityGroupFile(existingFile *entity.GroupFile, req *communityEvent.UpdateGroupFilePayload) *entity.GroupFile {
	hasChanges := false

	// Kiểm tra từng trường: Nếu khác nil (Client có gửi) VÀ khác với giá trị hiện tại
	if req.FileName != nil && *req.FileName != existingFile.FileName {
		existingFile.FileName = *req.FileName
		hasChanges = true
	}

	if req.Caption != nil && *req.Caption != existingFile.Caption {
		existingFile.Caption = *req.Caption
		hasChanges = true
	}

	if req.ThumbnailURL != nil && *req.ThumbnailURL != existingFile.ThumbnailURL {
		existingFile.ThumbnailURL = *req.ThumbnailURL
		hasChanges = true
	}

	// Best Practice: Chỉ cập nhật UpdatedAt nếu thực sự có dữ liệu bị thay đổi
	if hasChanges {
		existingFile.UpdatedAt = time.Now()
	}

	return existingFile
}
