package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ---------------------------------------------------------
// MAPPER REQ TO ENTITY
// ---------------------------------------------------------

// ToEntityGroupFile: Ánh xạ từ CreateGroupFileReq sang Entity
func ToEntityGroupFile(r *req.CreateGroupFileReq) *entity.GroupFile {
	if r == nil {
		return nil
	}

	groupID, _ := primitive.ObjectIDFromHex(r.GroupID) // Cần validate GroupID ở tầng Transport/Handler trước

	return &entity.GroupFile{
		ID:            primitive.NewObjectID(), // Tự sinh ID cho document mới
		GroupID:       groupID,
		UploaderID:    r.UploaderID,
		FileName:      r.FileName,
		FileType:      r.FileType,
		FileSize:      r.FileSize,
		StorageFileID: r.StorageFileID,
		DownloadCount: 0, // Mặc định khi mới tạo file thì số lượt tải là 0
		CreatedAt:     time.Now(),
	}
}

// UpdateToEntityGroupFile: Cập nhật đè các trường có trong UpdateGroupFileReq vào Entity hiện tại
func UpdateToEntityGroupFile(r *req.UpdateGroupFileReq, e *entity.GroupFile) {
	if r == nil || e == nil {
		return
	}

	if r.GroupID != nil {
		if oid, err := primitive.ObjectIDFromHex(*r.GroupID); err == nil {
			e.GroupID = oid
		}
	}
	if r.UploaderID != nil {
		e.UploaderID = *r.UploaderID
	}
	if r.FileName != nil {
		e.FileName = *r.FileName
	}
	if r.FileType != nil {
		e.FileType = *r.FileType
	}
	if r.FileSize != nil {
		e.FileSize = *r.FileSize
	}
	if r.StorageFileID != nil {
		e.StorageFileID = *r.StorageFileID
	}
	if r.DownloadCount != nil {
		e.DownloadCount = *r.DownloadCount
	}
}

// ---------------------------------------------------------
// MAPPER ENTITY TO RES
// ---------------------------------------------------------

// ToResGroupFile: Ánh xạ 100% từ Entity ra GroupFileRes
func ToResGroupFile(e *entity.GroupFile) *res.GroupFileRes {
	if e == nil {
		return nil
	}

	return &res.GroupFileRes{
		ID:            e.ID.Hex(),
		GroupID:       e.GroupID.Hex(),
		UploaderID:    e.UploaderID,
		FileName:      e.FileName,
		FileType:      e.FileType,
		FileSize:      e.FileSize,
		StorageFileID: e.StorageFileID,
		DownloadCount: e.DownloadCount,
		CreatedAt:     e.CreatedAt,
	}
}
