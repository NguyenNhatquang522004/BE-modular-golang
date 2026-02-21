package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/enum"
)

const (
	IndexSearchGroups = "search_groups"
)

// GeoLocation định nghĩa cấu trúc tọa độ chuẩn cho Elasticsearch và MongoDB
type GeoLocation struct {
	Lat float64 `json:"lat" bson:"lat"`
	Lon float64 `json:"lon" bson:"lon"`
}

// SearchGroup đại diện cho document trong index "search_groups"
type SearchGroup struct {
	// ID: Map với "_id" của ES.
	// omitempty: Để ES tự sinh ID nếu bạn không cung cấp, hoặc mapping từ Mongo ObjectID sang String.
	ID string `json:"_id,omitempty" bson:"_id,omitempty"`

	// Name: Tên nhóm
	Name string `json:"name" bson:"name"`

	// Description: Mô tả nhóm
	Description string `json:"description" bson:"description"`

	// Tags: Mảng các từ khóa (keyword/text)
	// Khởi tạo mảng rỗng thay vì nil nếu không có tag để tránh null trong JSON
	Tags []string `json:"tags" bson:"tags"`

	// Privacy: Enum quản lý quyền riêng tư.
	// Nhờ thư viện enumer, field này sẽ serialize thành string ("public")
	Privacy enum.PrivacySearch `json:"privacy" bson:"privacy"`

	// MemberCount: Số lượng thành viên
	MemberCount int `json:"member_count" bson:"member_count"`

	// Location: Tọa độ địa lý
	Location GeoLocation `json:"location" bson:"location"`

	// Metadata bổ sung (Optional - Best practice cho tracking)
	IndexedAt time.Time `json:"indexed_at,omitempty" bson:"-"` // Chỉ dùng cho ES, không lưu Mongo
}

func (SearchGroup) IndexName() string {
	return IndexSearchGroups
}
