package entity

import "time"

type ProfileES struct {
	// ID của ES nên trùng với ID của MongoDB (Hex String)
	ID     string `json:"id"`
	UserID string `json:"user_id"` // Keyword

	// Full-text search fields
	FullName  string `json:"full_name"`  // Text (Analyzer: standard/vietnamese)
	FirstName string `json:"first_name"` // Text
	LastName  string `json:"last_name"`  // Text
	Bio       string `json:"bio"`        // Text
	Slug      string `json:"slug"`       // Keyword

	// Filters
	DateOfBirth string `json:"date_of_birth,omitempty"`
	Gender      string     `json:"gender"` // Lưu String ("Male", "Female") thay vì Int để dễ filter/aggs

	// Media (Chỉ cần URL để hiển thị kết quả search, không cần nested object phức tạp)
	AvatarURL string `json:"avatar_url"`

	// Location (Geo-Search)
	AddressCity    string `json:"address_city"`    // Keyword
	AddressCountry string `json:"address_country"` // Keyword
	// ES Geo-point format: [lon, lat]
	Location []float64 `json:"location,omitempty"`

	// Searchable Lists (Nên flatten nếu chỉ cần search text đơn giản)
	// Hoặc giữ nguyên structure nếu dùng Nested Query
	Skills        []string        `json:"skills"`       // Ví dụ: Lấy từ CV hoặc Bio
	WorkHistory   []ESWorkHistory `json:"work_history"` // Nested Object
	EducationList []ESEducation   `json:"education_list"`

	// Meta
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Settings (Để filter user Private ra khỏi kết quả search)
	IsPrivate bool `json:"is_private"`
}

type ESWorkHistory struct {
	Company   string `json:"company"`  // Text
	Position  string `json:"position"` // Text
	IsCurrent bool   `json:"is_current"`
}

type ESEducation struct {
	Institution string `json:"institution"` // Text
	Degree      string `json:"degree"`
}
