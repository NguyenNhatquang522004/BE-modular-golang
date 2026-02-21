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

// ToEntityGroup: Ánh xạ 100% từ CreateGroupReq sang Entity
func ToEntityGroup(r *req.CreateGroupReq) *entity.Group {
	if r == nil {
		return nil
	}

	now := time.Now()
	e := &entity.Group{
		ID:          primitive.NewObjectID(), // Tự động generate ObjectID mới
		CreatorID:   r.CreatorID,
		Name:        r.Name,
		Slug:        r.Slug,
		Description: r.Description,
		Tags:        r.Tags,
		Cover: entity.GroupCover{
			URL:       r.Cover.URL,
			PositionY: r.Cover.PositionY,
		},
		Avatar: entity.GroupAvatar{
			URL: r.Avatar.URL,
		},
		Privacy:    r.Privacy,
		CategoryID: r.CategoryID,
		Settings: entity.GroupSettings{
			RequireApprovalToJoin: r.Settings.RequireApprovalToJoin,
			RequireApprovalToPost: r.Settings.RequireApprovalToPost,
			AllowMemberPosting:    r.Settings.AllowMemberPosting,
			WhoCanApproveMember:   r.Settings.WhoCanApproveMember,
		},
		CommunityChats: entity.FeatureFlag{
			IsEnabled: r.CommunityChats.IsEnabled,
		},
		MembershipQuestions: entity.FeatureFlag{
			IsEnabled: r.MembershipQuestions.IsEnabled,
		},
		Stats:     entity.GroupStats{}, // Mặc định tất cả các count = 0 khi mới tạo
		CreatedAt: now,
		UpdatedAt: now,
	}

	if len(r.Rules) > 0 {
		e.Rules = make([]entity.GroupRule, len(r.Rules))
		for i, rule := range r.Rules {
			e.Rules[i] = entity.GroupRule{
				Title:   rule.Title,
				Content: rule.Content,
			}
		}
	}

	return e
}

// UpdateToEntityGroup: Ánh xạ đè các trường có trong UpdateGroupReq vào Entity hiện tại (Partial Update)
func UpdateToEntityGroup(r *req.UpdateGroupReq, e *entity.Group) {
	if r == nil || e == nil {
		return
	}

	if r.Name != nil {
		e.Name = *r.Name
	}
	if r.Slug != nil {
		e.Slug = *r.Slug
	}
	if r.Description != nil {
		e.Description = *r.Description
	}
	if r.Tags != nil {
		e.Tags = *r.Tags
	}
	if r.Cover != nil {
		e.Cover.URL = r.Cover.URL
		e.Cover.PositionY = r.Cover.PositionY
	}
	if r.Avatar != nil {
		e.Avatar.URL = r.Avatar.URL
	}
	if r.Privacy != nil {
		e.Privacy = *r.Privacy
	}
	if r.CategoryID != nil {
		e.CategoryID = *r.CategoryID
	}
	if r.Rules != nil {
		e.Rules = make([]entity.GroupRule, len(*r.Rules))
		for i, rule := range *r.Rules {
			e.Rules[i] = entity.GroupRule{
				Title:   rule.Title,
				Content: rule.Content,
			}
		}
	}
	if r.Settings != nil {
		e.Settings.RequireApprovalToJoin = r.Settings.RequireApprovalToJoin
		e.Settings.RequireApprovalToPost = r.Settings.RequireApprovalToPost
		e.Settings.AllowMemberPosting = r.Settings.AllowMemberPosting
		e.Settings.WhoCanApproveMember = r.Settings.WhoCanApproveMember
	}
	if r.CommunityChats != nil {
		e.CommunityChats.IsEnabled = r.CommunityChats.IsEnabled
	}
	if r.MembershipQuestions != nil {
		e.MembershipQuestions.IsEnabled = r.MembershipQuestions.IsEnabled
	}

	e.UpdatedAt = time.Now() // Tự động cập nhật thời gian
}

// ---------------------------------------------------------
// MAPPER ENTITY TO RES
// ---------------------------------------------------------

// ToResGroup: Ánh xạ 100% từ Entity sang GroupRes
func ToResGroup(e *entity.Group) *res.GroupRes {
	if e == nil {
		return nil
	}

	r := &res.GroupRes{
		ID:          e.ID.Hex(), // Parse từ ObjectID sang String
		CreatorID:   e.CreatorID,
		Name:        e.Name,
		Slug:        e.Slug,
		Description: e.Description,
		Tags:        e.Tags,
		Cover: res.ResGroupCover{
			URL:       e.Cover.URL,
			PositionY: e.Cover.PositionY,
		},
		Avatar: res.ResGroupAvatar{
			URL: e.Avatar.URL,
		},
		Privacy:    e.Privacy,
		CategoryID: e.CategoryID,
		Settings: res.ResGroupSettings{
			RequireApprovalToJoin: e.Settings.RequireApprovalToJoin,
			RequireApprovalToPost: e.Settings.RequireApprovalToPost,
			AllowMemberPosting:    e.Settings.AllowMemberPosting,
			WhoCanApproveMember:   e.Settings.WhoCanApproveMember,
		},
		CommunityChats: res.ResFeatureFlag{
			IsEnabled: e.CommunityChats.IsEnabled,
		},
		MembershipQuestions: res.ResFeatureFlag{
			IsEnabled: e.MembershipQuestions.IsEnabled,
		},
		Stats: res.ResGroupStats{
			MemberCount:        e.Stats.MemberCount,
			PostCount:          e.Stats.PostCount,
			PendingMemberCount: e.Stats.PendingMemberCount,
			PendingPostCount:   e.Stats.PendingPostCount,
			ReportedPostCount:  e.Stats.ReportedPostCount,
		},
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		DeletedAt: e.DeletedAt,
	}

	if len(e.Rules) > 0 {
		r.Rules = make([]res.ResGroupRule, len(e.Rules))
		for i, rule := range e.Rules {
			r.Rules[i] = res.ResGroupRule{
				Title:   rule.Title,
				Content: rule.Content,
			}
		}
	}

	return r
}
