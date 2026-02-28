package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ========================
// MAPPER REQUEST -> ENTITY
// ========================

// ToEntityConversation tạo mới một Entity Conversation từ Request.
func ToEntityConversation(r *req.ConversationReq) *entity.Conversation {
	if r == nil {
		return nil
	}

	id, _ := primitive.ObjectIDFromHex(r.ID)
	if r.ID == "" {
		id = primitive.NewObjectID()
	}

	// Xử lý an toàn con trỏ RelatedGroupID
	var relatedGroupID *primitive.ObjectID
	if r.RelatedGroupID != nil && *r.RelatedGroupID != "" {
		oid, err := primitive.ObjectIDFromHex(*r.RelatedGroupID)
		if err == nil {
			relatedGroupID = &oid
		}
	}

	// Mappers cho các struct lồng nhau
	var avatar *entity.ConversationAvatar
	if r.Avatar != nil {
		avatar = &entity.ConversationAvatar{URL: r.Avatar.URL}
	}

	var theme *entity.ConversationTheme
	if r.Theme != nil {
		theme = &entity.ConversationTheme{
			Color:         r.Theme.Color,
			Emoji:         r.Theme.Emoji,
			BackgroundURL: r.Theme.BackgroundURL,
		}
	}

	var lastMessage *entity.LastMessageCache
	if r.LastMessage != nil {
		lastMessage = &entity.LastMessageCache{
			MessageID: r.LastMessage.MessageID,
			Content:   r.LastMessage.Content,
			SenderID:  r.LastMessage.SenderID,
			Type:      r.LastMessage.Type,
			CreatedAt: r.LastMessage.CreatedAt,
		}
	}

	now := time.Now()
	createdAt := r.CreatedAt
	if createdAt.IsZero() {
		createdAt = now
	}
	updatedAt := r.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = now
	}

	return &entity.Conversation{
		ID:             id,
		Type:           *r.Type,
		Scope:          *r.Scope,
		Status:         *r.Status,
		Name:           r.Name,
		Avatar:         avatar,
		CreatorID:      r.CreatorID,
		OwnerID:        r.OwnerID,
		RelatedGroupID: relatedGroupID,
		Permissions: entity.ConversationPermissions{
			SendMessage: *r.Permissions.SendMessage,
			AddMember:   *r.Permissions.AddMember,
		},
		Theme:            theme,
		LastMessage:      lastMessage,
		ParticipantCount: r.ParticipantCount,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}
}

// UpdateToEntityConversation cập nhật Entity hiện có từ Request.
func UpdateToEntityConversation(r *req.ConversationReq, e *entity.Conversation) {
	if r == nil || e == nil {
		return
	}

	if r.Type != nil {
		e.Type = *r.Type
	}
	if r.Scope != nil {
		e.Scope = *r.Scope
	}
	if r.Status != nil {
		e.Status = *r.Status
	}
	if r.Name != "" {
		e.Name = r.Name
	}

	if r.Avatar != nil {
		e.Avatar = &entity.ConversationAvatar{URL: r.Avatar.URL}
	}

	if r.CreatorID != "" {
		e.CreatorID = r.CreatorID
	}
	if r.OwnerID != "" {
		e.OwnerID = r.OwnerID
	}

	if r.RelatedGroupID != nil && *r.RelatedGroupID != "" {
		if oid, err := primitive.ObjectIDFromHex(*r.RelatedGroupID); err == nil {
			e.RelatedGroupID = &oid
		}
	}

	// Update Permissions
	if r.Permissions.SendMessage != nil {
		e.Permissions.SendMessage = *r.Permissions.SendMessage
	}
	if r.Permissions.AddMember != nil {
		e.Permissions.AddMember = *r.Permissions.AddMember
	}

	if r.Theme != nil {
		e.Theme = &entity.ConversationTheme{
			Color:         r.Theme.Color,
			Emoji:         r.Theme.Emoji,
			BackgroundURL: r.Theme.BackgroundURL,
		}
	}

	if r.LastMessage != nil {
		e.LastMessage = &entity.LastMessageCache{
			MessageID: r.LastMessage.MessageID,
			Content:   r.LastMessage.Content,
			SenderID:  r.LastMessage.SenderID,
			Type:      r.LastMessage.Type,
			CreatedAt: r.LastMessage.CreatedAt,
		}
	}

	if r.ParticipantCount > 0 { // Tuỳ logic có cho phép update = 0 hay không
		e.ParticipantCount = r.ParticipantCount
	}

	e.UpdatedAt = time.Now()
}

// ========================
// MAPPER -> RESPONSE
// ========================
// MAPPER -> RESPONSE
// ========================

// ToResFromReqConversation chuyển trực tiếp từ Req sang Res
func ToResFromReqConversation(r *req.ConversationReq) *res.ConversationRes {
	if r == nil {
		return nil
	}

	var avatarRes *res.ConversationAvatarRes
	if r.Avatar != nil {
		avatarRes = &res.ConversationAvatarRes{URL: r.Avatar.URL}
	}

	var themeRes *res.ConversationThemeRes
	if r.Theme != nil {
		themeRes = &res.ConversationThemeRes{
			Color:         r.Theme.Color,
			Emoji:         r.Theme.Emoji,
			BackgroundURL: r.Theme.BackgroundURL,
		}
	}

	var lastMessageRes *res.LastMessageCacheRes
	if r.LastMessage != nil {
		lastMessageRes = &res.LastMessageCacheRes{
			MessageID: r.LastMessage.MessageID,
			Content:   r.LastMessage.Content,
			SenderID:  r.LastMessage.SenderID,
			Type:      r.LastMessage.Type,
			CreatedAt: r.LastMessage.CreatedAt,
		}
	}

	return &res.ConversationRes{
		ID:             r.ID,
		Type:           *r.Type,
		Scope:          *r.Scope,
		Status:         *r.Status,
		Name:           r.Name,
		Avatar:         avatarRes,
		CreatorID:      r.CreatorID,
		OwnerID:        r.OwnerID,
		RelatedGroupID: r.RelatedGroupID,
		Permissions: res.ConversationPermissionsRes{
			SendMessage: *r.Permissions.SendMessage,
			AddMember:   *r.Permissions.AddMember,
		},
		Theme:            themeRes,
		LastMessage:      lastMessageRes,
		ParticipantCount: r.ParticipantCount,
		CreatedAt:        r.CreatedAt,
		UpdatedAt:        r.UpdatedAt,
	}
}

// ToResFromEntityConversation chuyển từ DB Entity sang Res để trả về Client.
func ToResFromEntityConversation(e *entity.Conversation) *res.ConversationRes {
	if e == nil {
		return nil
	}

	var relatedGroupStr *string
	if e.RelatedGroupID != nil {
		idHex := e.RelatedGroupID.Hex()
		relatedGroupStr = &idHex
	}

	var avatarRes *res.ConversationAvatarRes
	if e.Avatar != nil {
		avatarRes = &res.ConversationAvatarRes{URL: e.Avatar.URL}
	}

	var themeRes *res.ConversationThemeRes
	if e.Theme != nil {
		themeRes = &res.ConversationThemeRes{
			Color:         e.Theme.Color,
			Emoji:         e.Theme.Emoji,
			BackgroundURL: e.Theme.BackgroundURL,
		}
	}

	var lastMessageRes *res.LastMessageCacheRes
	if e.LastMessage != nil {
		lastMessageRes = &res.LastMessageCacheRes{
			MessageID: e.LastMessage.MessageID,
			Content:   e.LastMessage.Content,
			SenderID:  e.LastMessage.SenderID,
			Type:      e.LastMessage.Type,
			CreatedAt: e.LastMessage.CreatedAt,
		}
	}

	return &res.ConversationRes{
		ID:             e.ID.Hex(),
		Type:           e.Type,
		Scope:          e.Scope,
		Status:         e.Status,
		Name:           e.Name,
		Avatar:         avatarRes,
		CreatorID:      e.CreatorID,
		OwnerID:        e.OwnerID,
		RelatedGroupID: relatedGroupStr,
		Permissions: res.ConversationPermissionsRes{
			SendMessage: e.Permissions.SendMessage,
			AddMember:   e.Permissions.AddMember,
		},
		Theme:            themeRes,
		LastMessage:      lastMessageRes,
		ParticipantCount: e.ParticipantCount,
		CreatedAt:        e.CreatedAt,
		UpdatedAt:        e.UpdatedAt,
	}
}
