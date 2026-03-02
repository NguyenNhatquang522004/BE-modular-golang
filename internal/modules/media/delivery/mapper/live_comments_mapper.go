package mapper

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
)

// ToEntity ánh xạ từ DTO Request sang Entity Cassandra
func ToEntityLiveComment(r *req.LiveCommentReq) *entity.LiveComment {
	if r == nil {
		return nil
	}

	// Convert []enum.UserBadge -> []string cho Cassandra

	return &entity.LiveComment{
		StreamID:      r.StreamID,
		CreatedAt:     r.CreatedAt,
		CommentID:     r.CommentID,
		UserID:        r.UserID,
		UserNickname:  r.UserNickname,
		UserAvatarURL: r.UserAvatarURL,
		UserBadges:    r.UserBadges,
		Content:       r.Content,
		IsPinned:      r.IsPinned,
	}
}

// UpdateToEntity cập nhật dữ liệu từ DTO Request vào Entity có sẵn
func UpdateToEntityLiveComment(r *req.LiveCommentReq, ent *entity.LiveComment) {
	if r == nil || ent == nil {
		return
	}
	ent.StreamID = r.StreamID
	ent.CreatedAt = r.CreatedAt
	ent.CommentID = r.CommentID
	ent.UserID = r.UserID
	ent.UserNickname = r.UserNickname
	ent.UserAvatarURL = r.UserAvatarURL
	ent.UserBadges = r.UserBadges
	ent.Content = r.Content
	ent.IsPinned = r.IsPinned
}
