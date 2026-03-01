// File: mapper/message_mapper.go
package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/entity"
)

// =========================================================================
// MAPPER REQ
// =========================================================================
func GenerateBucket(t time.Time) int {
	// Ép về UTC để đảm bảo tính nhất quán của dữ liệu trên toàn hệ thống
	tUTC := t.UTC()
	return tUTC.Year()*100 + int(tUTC.Month())
}

// ToEntityMessage chuyển hoàn toàn từ Req sang Entity.
func ToEntityMessage(r *req.MessageReq, conversationID string) *entity.Message {
	if r == nil {
		return nil
	}
	now := time.Now().UTC()
	return &entity.Message{
		ConversationID:   conversationID,
		Bucket:           GenerateBucket(now),
		MessageID:        r.MessageID,
		SenderID:         r.SenderID,
		Type:             r.Type,
		Content:          r.Content,
		Attachments:      r.Attachments,
		IsEdited:         r.IsEdited,
		ReplyToMessageID: r.ReplyToMessageID,
		StoryRefID:       r.StoryRefID,
		IsRevoked:        r.IsRevoked,
		CreatedAt:        r.CreatedAt,
	}
}

// UpdateToEntityMessage map dữ liệu từ Req vào một Entity đã có sẵn.
func UpdateToEntityMessage(r *req.MessageReq, e *entity.Message) {
	if r == nil || e == nil {
		return
	}
	e.ConversationID = r.ConversationID
	e.Bucket = r.Bucket
	e.MessageID = r.MessageID
	e.SenderID = r.SenderID
	e.Type = r.Type
	e.Content = r.Content
	e.Attachments = r.Attachments
	e.IsEdited = r.IsEdited
	e.ReplyToMessageID = r.ReplyToMessageID
	e.StoryRefID = r.StoryRefID
	e.IsRevoked = r.IsRevoked
	e.CreatedAt = r.CreatedAt
}

// =========================================================================
// MAPPER RES
// =========================================================================

// ReqToRes map trực tiếp từ Request sang Response (Theo đúng yêu cầu của bạn).
func ToResFromReqMessage(r *req.MessageReq) *res.MessageRes {
	if r == nil {
		return nil
	}
	return &res.MessageRes{
		ConversationID:   r.ConversationID,
		Bucket:           r.Bucket,
		MessageID:        r.MessageID,
		SenderID:         r.SenderID,
		Type:             r.Type,
		Content:          r.Content,
		Attachments:      r.Attachments,
		IsEdited:         r.IsEdited,
		ReplyToMessageID: r.ReplyToMessageID,
		StoryRefID:       r.StoryRefID,
		IsRevoked:        r.IsRevoked,
		CreatedAt:        r.CreatedAt,
	}
}

// EntityToResMessage map từ Entity trong DB ra Response DTO (Hàm cần thiết khi GET từ Cassandra).
func EntityToResMessage(e *entity.Message) *res.MessageRes {
	if e == nil {
		return nil
	}
	return &res.MessageRes{
		ConversationID:   e.ConversationID,
		Bucket:           e.Bucket,
		MessageID:        e.MessageID,
		SenderID:         e.SenderID,
		Type:             e.Type,
		Content:          e.Content,
		Attachments:      e.Attachments,
		IsEdited:         e.IsEdited,
		ReplyToMessageID: e.ReplyToMessageID,
		StoryRefID:       e.StoryRefID,
		IsRevoked:        e.IsRevoked,
		CreatedAt:        e.CreatedAt,
	}
}
