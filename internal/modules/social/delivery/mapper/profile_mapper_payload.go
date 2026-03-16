package mapper

import (
	"regexp"
	"strings"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/socialEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToEntityProfilePayload(req *socialEvent.ProfilePayload) (*entity.Profiles, error) {
	now := time.Now()

	// 1. Tạo ID và xử lý dữ liệu dẫn xuất
	profileID := primitive.NewObjectID()
	fullName := strings.TrimSpace(req.FirstName + " " + req.LastName)
	slug := generateSlug(fullName)

	// 2. Map các trường cơ bản
	profile := &entity.Profiles{
		ID:          profileID,
		UserID:      req.UserID,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		FullName:    fullName,
		Slug:        slug,
		Bio:         req.Bio,
		DateOfBirth: req.DateOfBirth,
		Gender:      req.Gender,
		CreatedAt:   now,
		UpdatedAt:   now,
		// Default Settings
		Settings: entity.ProfileSettings{
			IsPrivate:         false,
			AllowSearchEngine: true,
		},
	}

	// 3. Map Sub-structs (Single Objects)
	// Sử dụng helper function để code gọn và an toàn với null
	profile.Avatar = mapAvatarPayload(req.Avatar)
	profile.CoverPhoto = mapCoverPhotoPayload(req.CoverPhoto)
	profile.Address = mapAddressPayload(req.Address)
	profile.SocialLinks = mapSocialLinksPayload(req.SocialLinks)
	profile.PhoneNumber = req.PhoneNumber

	// 4. Map CV Document (Cần validate ID hex string)
	if req.CVDocument != nil {
		cvDoc, err := mapCVDocumentPayload(req.CVDocument)
		if err != nil {
			return nil, err // Trả về lỗi nếu FileID không hợp lệ
		}
		profile.CVDocument = cvDoc
	}

	// 5. Map Arrays (WorkExperience & Education)
	// Vì là Create, ta luôn tạo ID mới cho các item trong mảng
	profile.WorkExperience = mapWorkExperienceCreatePayload(req.WorkExperience)
	profile.Education = mapEducationCreatePayload(req.Education)

	// 6. Map Settings & Notifications (Override default nếu có)
	if req.Settings != nil {
		profile.Settings.IsPrivate = req.Settings.IsPrivate
		profile.Settings.AllowSearchEngine = req.Settings.AllowSearchEngine
	}

	return profile, nil
}

// ToEntityUpdateProfilePayload: Cập nhật Entity hiện tại từ DTO (Update)
func ToEntityUpdateProfilePayload(current *entity.Profiles, req *socialEvent.ProfilePayload) (*entity.Profiles, error) {
	now := time.Now()

	// 1. Cập nhật thông tin cơ bản
	current.FirstName = req.FirstName
	current.LastName = req.LastName
	current.Bio = req.Bio
	current.DateOfBirth = req.DateOfBirth
	current.Gender = req.Gender
	current.PhoneNumber = req.PhoneNumber
	current.UpdatedAt = now

	// 2. Logic cập nhật FullName & Slug
	// Chỉ tính toán lại khi tên thay đổi để tối ưu hiệu năng
	newFullName := strings.TrimSpace(req.FirstName + " " + req.LastName)
	if current.FullName != newFullName {
		current.FullName = newFullName
		current.Slug = generateSlug(newFullName)
	}

	// 3. Cập nhật Sub-structs (Pointer)
	// Logic: Nếu DTO gửi lên != nil thì thay thế. Nếu nil thì giữ nguyên cái cũ.
	// (Hoặc tùy logic business, nếu muốn xóa thì FE cần gửi object rỗng, ở đây ta assume là replace if present)
	if req.Avatar != nil {
		current.Avatar = mapAvatarPayload(req.Avatar)
	}
	if req.CoverPhoto != nil {
		current.CoverPhoto = mapCoverPhotoPayload(req.CoverPhoto)
	}
	if req.Address != nil {
		current.Address = mapAddressPayload(req.Address)
	}
	if req.SocialLinks != nil {
		current.SocialLinks = mapSocialLinksPayload(req.SocialLinks)
	}

	// 4. Cập nhật CV
	if req.CVDocument != nil {
		cvDoc, err := mapCVDocumentPayload(req.CVDocument)
		if err != nil {
			return nil, err
		}
		current.CVDocument = cvDoc
	}

	// 5. Cập nhật Mảng (SMART UPDATE STRATEGY)
	// Giữ nguyên ID của các item cũ, tạo ID mới cho item mới, xóa item không có trong list
	if req.WorkExperience != nil {
		current.WorkExperience = mapWorkExperienceSmartPayload(current.WorkExperience, req.WorkExperience)
	}
	if req.Education != nil {
		current.Education = mapEducationSmartPayload(current.Education, req.Education)
	}

	// 6. Cập nhật Settings
	if req.Settings != nil {
		current.Settings = entity.ProfileSettings{
			IsPrivate:         req.Settings.IsPrivate,
			AllowSearchEngine: req.Settings.AllowSearchEngine,
		}
	}

	return current, nil
}

// =================================================================================
// PRIVATE HELPER MAPPERS (Sub-structs)
// =================================================================================

func mapAvatarPayload(req *socialEvent.AvatarPayload) *entity.ProfileAvatar {
	if req == nil {
		return nil
	}
	return &entity.ProfileAvatar{
		ID:        primitive.NewObjectID(), // Tạo ID mới cho avatar (liên kết với GridFS)
		URL:       req.URL,
		UpdatedAt: time.Now(),
	}
}
func mapCoverPhotoPayload(req *socialEvent.CoverPhotoPayload) *entity.ProfileCover {
	if req == nil {
		return nil
	}
	return &entity.ProfileCover{
		ID:        primitive.NewObjectID(), // Tạo ID mới cho cover photo (liên kết với GridFS)
		URL:       req.URL,
		PositionY: req.PositionY,
	}
}
func mapAddressPayload(req *socialEvent.AddressPayload) *entity.ProfileAddress {
	if req == nil {
		return nil
	}
	return &entity.ProfileAddress{
		Street:      req.Street,
		City:        req.City,
		Country:     req.Country,
		Coordinates: req.Coordinates,
	}
}

func mapSocialLinksPayload(req *socialEvent.SocialLinksPayload) *entity.SocialLinks {
	if req == nil {
		return nil
	}
	return &entity.SocialLinks{
		Facebook:  req.Facebook,
		Twitter:   req.Twitter,
		Linkedin:  req.Linkedin,
		Instagram: req.Instagram,
		Github:    req.Github,
		Website:   req.Website,
	}
}

func mapCVDocumentPayload(req *socialEvent.CVDocumentPayload) (*entity.CVDocument, error) {
	if req == nil {
		return nil, nil
	}
	oid, err := primitive.ObjectIDFromHex(req.FileID)
	if err != nil {
		return nil, err
	}
	return &entity.CVDocument{
		FileID:     oid,
		Filename:   req.Filename,
		UploadedAt: time.Now(),
	}, nil
}

// =================================================================================
// ARRAY MAPPERS (Create & Smart Update)
// =================================================================================

// mapWorkExperienceCreatePayload: Dùng cho tạo mới Profile
func mapWorkExperienceCreatePayload(reqList []*socialEvent.WorkExperiencePayload) []*entity.WorkExperience {
	if len(reqList) == 0 {
		return []*entity.WorkExperience{} // Tránh trả về nil cho slice
	}
	result := make([]*entity.WorkExperience, len(reqList))
	for i, item := range reqList {
		result[i] = &entity.WorkExperience{
			ID:          primitive.NewObjectID(), // Luôn tạo ID mới
			Company:     item.Company,
			Position:    item.Position,
			StartDate:   item.StartDate,
			EndDate:     item.EndDate,
			IsCurrent:   item.IsCurrent,
			Description: item.Description,
		}
	}
	return result
}

// mapWorkExperienceSmartPayload: Dùng cho Update Profile (Giữ ID cũ)
func mapWorkExperienceSmartPayload(currentList []*entity.WorkExperience, reqList []*socialEvent.WorkExperiencePayload) []*entity.WorkExperience {
	if len(reqList) == 0 {
		return []*entity.WorkExperience{}
	}

	// 1. Indexing existing items by ID for O(1) lookup
	currentMap := make(map[string]*entity.WorkExperience)
	for _, item := range currentList {
		currentMap[item.ID.Hex()] = item
	}

	var result []*entity.WorkExperience
	for _, reqItem := range reqList {
		var itemToSave *entity.WorkExperience

		// 2. Check update existing
		if reqItem.ID != "" {
			if existing, exists := currentMap[reqItem.ID]; exists {
				existing.Company = reqItem.Company
				existing.Position = reqItem.Position
				existing.StartDate = reqItem.StartDate
				existing.EndDate = reqItem.EndDate
				existing.IsCurrent = reqItem.IsCurrent
				existing.Description = reqItem.Description
				itemToSave = existing
			}
		}

		// 3. Create new if ID mismatch or empty
		if itemToSave == nil {
			itemToSave = &entity.WorkExperience{
				ID:          primitive.NewObjectID(),
				Company:     reqItem.Company,
				Position:    reqItem.Position,
				StartDate:   reqItem.StartDate,
				EndDate:     reqItem.EndDate,
				IsCurrent:   reqItem.IsCurrent,
				Description: reqItem.Description,
			}
		}
		result = append(result, itemToSave)
	}
	return result
}

// Tương tự cho Education
func mapEducationCreatePayload(reqList []*socialEvent.EducationPayload) []*entity.Education {
	if len(reqList) == 0 {
		return []*entity.Education{}
	}
	result := make([]*entity.Education, len(reqList))
	for i, item := range reqList {
		result[i] = &entity.Education{
			ID:           primitive.NewObjectID(),
			Institution:  item.Institution,
			Degree:       item.Degree,
			FieldOfStudy: item.FieldOfStudy,
			StartDate:    item.StartDate,
			EndDate:      item.EndDate,
		}
	}
	return result
}

func mapEducationSmartPayload(currentList []*entity.Education, reqList []*socialEvent.EducationPayload) []*entity.Education {
	if len(reqList) == 0 {
		return []*entity.Education{}
	}

	currentMap := make(map[string]*entity.Education)
	for _, item := range currentList {
		currentMap[item.ID.Hex()] = item
	}

	var result []*entity.Education
	for _, reqItem := range reqList {
		var itemToSave *entity.Education

		if reqItem.ID != "" {
			if existing, exists := currentMap[reqItem.ID]; exists {
				existing.Institution = reqItem.Institution
				existing.Degree = reqItem.Degree
				existing.FieldOfStudy = reqItem.FieldOfStudy
				existing.StartDate = reqItem.StartDate
				existing.EndDate = reqItem.EndDate
				itemToSave = existing
			}
		}

		if itemToSave == nil {
			itemToSave = &entity.Education{
				ID:           primitive.NewObjectID(),
				Institution:  reqItem.Institution,
				Degree:       reqItem.Degree,
				FieldOfStudy: reqItem.FieldOfStudy,
				StartDate:    reqItem.StartDate,
				EndDate:      reqItem.EndDate,
			}
		}
		result = append(result, itemToSave)
	}
	return result
}

// =================================================================================
// UTILS
// =================================================================================

// generateSlugPayload tạo URL thân thiện.
// Nên dùng thư viện "github.com/gosimple/slug" để xử lý tiếng Việt tốt hơn.
func generateSlugPayload(input string) string {
	// Simple implementation
	s := strings.ToLower(input)
	// Replace vietnamese chars (basic example)
	s = strings.ReplaceAll(s, "đ", "d")
	// Remove non-alphanumeric
	reg, _ := regexp.Compile("[^a-z0-9 ]+")
	s = reg.ReplaceAllString(s, "")
	// Replace spaces with hyphens
	s = strings.ReplaceAll(strings.TrimSpace(s), " ", "-")
	return s
}
