package mapper

import (
	"strings"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/communityEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ToCreateGroupEntity chuyển đổi CreateGroupPayload thành entity.Group
// Truyền thêm creatorID lấy từ Token/Context để đảm bảo bảo mật.
func ToCreateGroupEntity(req *communityEvent.CreateGroupPayload, creatorID string) *entity.Group {
	now := time.Now()
	if req.GroupID == "" {
		return nil // Hoặc trả về lỗi nếu GroupID là bắt buộc
	}
	convertGroupID, err := primitive.ObjectIDFromHex(req.GroupID)
	if err != nil {
		return nil // Hoặc trả về lỗi nếu GroupID không hợp lệ
	}
	group := &entity.Group{
		ID:         convertGroupID,
		CreatorID:  creatorID,
		Name:       req.Name,
		Slug:       generateSlug(req.Name), // Helper func để tạo URL friendly
		Privacy:    req.Privacy,
		CategoryID: req.CategoryID,
		// Xử lý slice an toàn, tránh nil slice trong DB nếu cần
		Tags:        safeStringSlice(req.Tags),
		Description: req.Description,
		Stats: entity.GroupStats{
			MemberCount:        1, // Creator auto join
			PostCount:          0,
			PendingMemberCount: 0,
			PendingPostCount:   0,
			ReportedPostCount:  0,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 1. Map Media (Cover, Avatar)
	if req.Cover != nil {
		group.Cover = entity.GroupCover{
			URL:       req.Cover.URL,
			PositionY: req.Cover.PositionY,
		}
	}
	if req.Avatar != nil {
		group.Avatar = entity.GroupAvatar{
			URL: req.Avatar.URL,
		}
	}

	// 2. Map Rules
	if len(req.Rules) > 0 {
		rules := make([]entity.GroupRule, 0, len(req.Rules))
		for _, r := range req.Rules {
			rules = append(rules, entity.GroupRule{
				Title:   r.Title,
				Content: r.Content,
			})
		}
		group.Rules = rules
	} else {
		group.Rules = []entity.GroupRule{}
	}

	// 3. Map Settings
	if req.Settings != nil {
		group.Settings = entity.GroupSettings{
			RequireApprovalToJoin: req.Settings.RequireApprovalToJoin,
			RequireApprovalToPost: req.Settings.RequireApprovalToPost,
			AllowMemberPosting:    req.Settings.AllowMemberPosting,
			WhoCanApproveMember:   req.Settings.WhoCanApproveMember,
		}
	}

	// 4. Map Feature Flags
	if req.CommunityChats != nil {
		group.CommunityChats = entity.FeatureFlag{IsEnabled: req.CommunityChats.IsEnabled}
	}
	if req.MembershipQuestions != nil {
		group.MembershipQuestions = entity.FeatureFlag{IsEnabled: req.MembershipQuestions.IsEnabled}
	}

	return group
}

// MapUpdatePayloadToEntity cập nhật các trường được truyền lên từ payload vào entity gốc.
// Sử dụng pattern Pass-by-Pointer để modify trực tiếp entity.
func MapUpdatePayloadToEntity(req *communityEvent.UpdateGroupPayload, group *entity.Group) {
	// Cập nhật Timestamps luôn luôn
	group.UpdatedAt = time.Now()

	// 1. Basic Fields
	if req.Name != nil {
		group.Name = *req.Name
		group.Slug = generateSlug(*req.Name) // Cập nhật lại Slug nếu tên đổi
	}
	if req.Description != nil {
		group.Description = *req.Description
	}
	if req.Tags != nil {
		group.Tags = *req.Tags
	}
	if req.Privacy != nil {
		group.Privacy = *req.Privacy
	}
	if req.CategoryID != nil {
		group.CategoryID = *req.CategoryID
	}

	// 2. Media (Kiểm tra nested pointers)
	if req.Cover != nil {
		if req.Cover.URL != nil {
			group.Cover.URL = *req.Cover.URL
		}
		if req.Cover.PositionY != nil {
			group.Cover.PositionY = *req.Cover.PositionY
		}
	}

	if req.Avatar != nil {
		if req.Avatar.URL != nil {
			group.Avatar.URL = *req.Avatar.URL
		}
	}

	// 3. Rules (PATCH chuẩn thường là Replace toàn bộ mảng nếu được truyền lên)
	if req.Rules != nil {
		rules := make([]entity.GroupRule, 0, len(*req.Rules))
		for _, r := range *req.Rules {
			rule := entity.GroupRule{}
			if r.Title != nil {
				rule.Title = *r.Title
			}
			if r.Content != nil {
				rule.Content = *r.Content
			}
			rules = append(rules, rule)
		}
		group.Rules = rules
	}

	// 4. Settings
	if req.Settings != nil {
		if req.Settings.RequireApprovalToJoin != nil {
			group.Settings.RequireApprovalToJoin = *req.Settings.RequireApprovalToJoin
		}
		if req.Settings.RequireApprovalToPost != nil {
			group.Settings.RequireApprovalToPost = *req.Settings.RequireApprovalToPost
		}
		if req.Settings.AllowMemberPosting != nil {
			group.Settings.AllowMemberPosting = *req.Settings.AllowMemberPosting
		}
		if req.Settings.WhoCanApproveMember != nil {
			group.Settings.WhoCanApproveMember = *req.Settings.WhoCanApproveMember
		}
	}

	// 5. Feature Flags
	if req.CommunityChats != nil {
		if req.CommunityChats.IsEnabled != nil {
			group.CommunityChats.IsEnabled = *req.CommunityChats.IsEnabled
		}
	}
	if req.MembershipQuestions != nil {
		if req.MembershipQuestions.IsEnabled != nil {
			group.MembershipQuestions.IsEnabled = *req.MembershipQuestions.IsEnabled
		}
	}
}

// --- HELPER FUNCTIONS ---

// generateSlug là hàm dummy, bạn nên dùng thư viện như github.com/gosimple/slug
func generateSlug(name string) string {
	// Giả lập cơ bản: "Học Golang" -> "hoc-golang-timestamp" để đảm bảo unique
	s := strings.ToLower(name)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

// safeStringSlice đảm bảo trả về mảng rỗng thay vì nil
func safeStringSlice(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}
