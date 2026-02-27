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
	var badges []string
	if len(r.UserBadges) > 0 {
		badges = make([]string, len(r.UserBadges))
		for i, b := range r.UserBadges {
			badges[i] = string(*b)
		}
	} else {
		badges = []string{} // Khởi tạo slice rỗng thay vì nil
	}

	return &entity.LiveComment{
		StreamID:      r.StreamID,
		CreatedAt:     r.CreatedAt,
		CommentID:     r.CommentID,
		UserID:        r.UserID,
		UserNickname:  r.UserNickname,
		UserAvatarURL: r.UserAvatarURL,
		UserBadges:    badges,
		Content:       r.Content,
		IsPinned:      r.IsPinned,
	}
}

// UpdateToEntity cập nhật dữ liệu từ DTO Request vào Entity có sẵn
func UpdateToEntityLiveComment(r *req.LiveCommentReq, ent *entity.LiveComment) {
	if r == nil || ent == nil {
		return
	}

	var badges []string
	if len(r.UserBadges) > 0 {
		badges = make([]string, len(r.UserBadges))
		for i, b := range r.UserBadges {
			badges[i] = string(*b)
		}
	} else {
		badges = []string{}
	}

	ent.StreamID = r.StreamID
	ent.CreatedAt = r.CreatedAt
	ent.CommentID = r.CommentID
	ent.UserID = r.UserID
	ent.UserNickname = r.UserNickname
	ent.UserAvatarURL = r.UserAvatarURL
	ent.UserBadges = badges
	ent.Content = r.Content
	ent.IsPinned = r.IsPinned
}
