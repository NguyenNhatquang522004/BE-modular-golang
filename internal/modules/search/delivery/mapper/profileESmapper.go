package mapper

import (
	entityES "github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/search/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
)

func ToElasticProfile(m *entity.Profiles) *entityES.ProfileES {
	if m == nil {
		return nil
	}

	// 1. Convert Enum Gender -> String
	genderStr := "Other"
	switch m.Gender {
	case enum.GenderMale:
		genderStr = "Male"
	case enum.GenderFemale:
		genderStr = "Female"
	case enum.GenderHidden:
		genderStr = "Hidden"
	}

	// 2. Handle Avatar
	avatarURL := ""
	if m.Avatar != nil {
		avatarURL = m.Avatar.URL
	}

	// 3. Handle Address & Location
	var city, country string
	var location []float64
	if m.Address != nil {
		city = m.Address.City
		country = m.Address.Country
		// Mongo: [Long, Lat] -> ES: [Long, Lat]. Thứ tự giống nhau nên gán trực tiếp được.
		// Lưu ý: Validate coordinates valid range (-180/180)
		if len(m.Address.Coordinates) == 2 {
			location = m.Address.Coordinates
		}
	}
	
	// 4. Map Work Experience
	var workHistory []entityES.ESWorkHistory
	for _, w := range m.WorkExperience {
		if w != nil {
			workHistory = append(workHistory, entityES.ESWorkHistory{
				Company:   w.Company,
				Position:  w.Position,
				IsCurrent: w.IsCurrent,
			})
		}
	}

	// 5. Map Education
	var eduList []entityES.ESEducation
	for _, e := range m.Education {
		if e != nil {
			eduList = append(eduList, entityES.ESEducation{
				Institution: e.Institution,
				Degree:      e.Degree,
			})
		}
	}

	// 6. Return ES Entity
	return &entityES.ProfileES{
		ID:             m.ID.Hex(), // Convert ObjectID to String
		UserID:         m.UserID,
		FullName:       m.FullName,
		FirstName:      m.FirstName,
		LastName:       m.LastName,
		Bio:            m.Bio,
		Slug:           m.Slug,
		DateOfBirth:    m.DateOfBirth,
		Gender:         genderStr,
		AvatarURL:      avatarURL,
		AddressCity:    city,
		AddressCountry: country,
		Location:       location,
		WorkHistory:    workHistory,
		EducationList:  eduList,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
		IsPrivate:      m.Settings.IsPrivate,
	}
}
