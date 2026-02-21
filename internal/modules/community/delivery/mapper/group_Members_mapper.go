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

// ToEntityGroupMember: Ánh xạ 100% từ CreateGroupMemberReq sang Entity GroupMember
func ToEntityGroupMember(r *req.CreateGroupMemberReq) *entity.GroupMember {
	if r == nil {
		return nil
	}

	// Parse String sang ObjectID (Bỏ qua lỗi trong mapper, nên handle lỗi ở tầng validation/handler trước)
	groupID, _ := primitive.ObjectIDFromHex(r.GroupID)

	now := time.Now()
	e := &entity.GroupMember{
		ID:           primitive.NewObjectID(),
		GroupID:      groupID,
		UserID:       r.UserID,
		Role:         r.Role,
		Status:       r.Status,
		InviterID:    r.InviterID,
		Badges:       r.Badges,
		JoinedAt:     now,
		LastActiveAt: now,
	}

	// Ánh xạ mảng JoinAnswers
	if len(r.JoinAnswers) > 0 {
		e.JoinAnswers = make([]entity.JoinAnswer, len(r.JoinAnswers))
		for i, ans := range r.JoinAnswers {
			qID, _ := primitive.ObjectIDFromHex(ans.QuestionID)
			e.JoinAnswers[i] = entity.JoinAnswer{
				QuestionID: qID,
				Answer:     ans.Answer,
			}
		}
	}

	// Ánh xạ DisciplineInfo
	if r.DisciplineInfo != nil {
		e.DisciplineInfo = &entity.DisciplineInfo{
			Reason:    r.DisciplineInfo.Reason,
			BannedBy:  r.DisciplineInfo.BannedBy,
			UntilDate: r.DisciplineInfo.UntilDate,
		}
	}

	return e
}

// UpdateToEntityGroupMember: Cập nhật đè (Partial Update) từ UpdateGroupMemberReq vào Entity có sẵn
func UpdateToEntityGroupMember(r *req.UpdateGroupMemberReq, e *entity.GroupMember) {
	if r == nil || e == nil {
		return
	}

	if r.GroupID != nil {
		if oid, err := primitive.ObjectIDFromHex(*r.GroupID); err == nil {
			e.GroupID = oid
		}
	}
	if r.UserID != nil {
		e.UserID = *r.UserID
	}
	if r.Role != nil {
		e.Role = *r.Role
	}
	if r.Status != nil {
		e.Status = *r.Status
	}
	if r.InviterID != nil {
		e.InviterID = *r.InviterID
	}
	if r.Badges != nil {
		e.Badges = *r.Badges
	}

	if r.JoinAnswers != nil {
		e.JoinAnswers = make([]entity.JoinAnswer, len(*r.JoinAnswers))
		for i, ans := range *r.JoinAnswers {
			qID, _ := primitive.ObjectIDFromHex(ans.QuestionID)
			e.JoinAnswers[i] = entity.JoinAnswer{
				QuestionID: qID,
				Answer:     ans.Answer,
			}
		}
	}

	// Ghi đè thông tin kỷ luật (nếu được truyền lên)
	if r.DisciplineInfo != nil {
		e.DisciplineInfo = &entity.DisciplineInfo{
			Reason:    r.DisciplineInfo.Reason,
			BannedBy:  r.DisciplineInfo.BannedBy,
			UntilDate: r.DisciplineInfo.UntilDate,
		}
	}
}

// ---------------------------------------------------------
// MAPPER ENTITY TO RES
// ---------------------------------------------------------

// ToResGroupMember: Ánh xạ 100% từ Entity sang DTO Response (Chuyển ObjectID thành String)
func ToResGroupMember(e *entity.GroupMember) *res.GroupMemberRes {
	if e == nil {
		return nil
	}

	response := &res.GroupMemberRes{
		ID:           e.ID.Hex(),
		GroupID:      e.GroupID.Hex(),
		UserID:       e.UserID,
		Role:         e.Role,
		Status:       e.Status,
		InviterID:    e.InviterID,
		Badges:       e.Badges,
		JoinedAt:     e.JoinedAt,
		LastActiveAt: e.LastActiveAt,
	}

	// Ánh xạ ngược JoinAnswers (Từ ObjectID sang String)
	if len(e.JoinAnswers) > 0 {
		response.JoinAnswers = make([]res.ResJoinAnswer, len(e.JoinAnswers))
		for i, ans := range e.JoinAnswers {
			response.JoinAnswers[i] = res.ResJoinAnswer{
				QuestionID: ans.QuestionID.Hex(),
				Answer:     ans.Answer,
			}
		}
	}

	if e.DisciplineInfo != nil {
		response.DisciplineInfo = &res.ResDisciplineInfo{
			Reason:    e.DisciplineInfo.Reason,
			BannedBy:  e.DisciplineInfo.BannedBy,
			UntilDate: e.DisciplineInfo.UntilDate,
		}
	}

	return response
}
