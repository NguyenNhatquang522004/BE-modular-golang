package mapper

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
	"github.com/gocql/gocql"
)

func ToEntityLiveCommentPayload(r *mediaEvent.LiveCommentPayload) *entity.LiveComment {
	var entitylivecomment = &entity.LiveComment{}
	if r == nil {
		return nil
	}
	convertSteamID, err := gocql.ParseUUID(r.StreamID)
	if err != nil {
		return nil
	}
	convertUserid, err := gocql.ParseUUID(r.UserID)
	if err != nil {
		return nil
	}
	if r.CommentID != "" {
		commentID, err := gocql.ParseUUID(r.CommentID)
		if err != nil {
			return nil
		}
		entitylivecomment.CommentID = commentID
	} else {
		entitylivecomment.CommentID = gocql.TimeUUID()
	}
	entitylivecomment.StreamID = convertSteamID
	entitylivecomment.CreatedAt = r.CreatedAt
	entitylivecomment.UserID = convertUserid
	entitylivecomment.UserNickname = r.UserNickname
	entitylivecomment.UserAvatarURL = r.UserAvatarURL
	entitylivecomment.UserBadges = r.UserBadges
	entitylivecomment.Content = r.Content
	entitylivecomment.IsPinned = r.IsPinned
	return entitylivecomment
}
func UpdateToEntityLiveCommentPayload(r *mediaEvent.LiveCommentPayload, ent *entity.LiveComment) {
	if r == nil || ent == nil {
		return
	}
	ent.UserNickname = r.UserNickname
	ent.UserAvatarURL = r.UserAvatarURL
	ent.UserBadges = r.UserBadges
	ent.Content = r.Content
	ent.IsPinned = r.IsPinned
}
