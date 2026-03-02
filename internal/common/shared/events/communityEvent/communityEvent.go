package communityEvent

import "github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"

type GroupStatsPayload struct {
	GroupID            string              `json:"group_id"`
	MemberCount        int                 `json:"member_count"`
	PostCount          int                 `json:"post_count"`
	PendingMemberCount int                 `json:"pending_member_count"`
	PendingPostCount   int                 `json:"pending_post_count"`
	ReportedPostCount  int                 `json:"reported_post_count"`
	EventType          constants.EventType `json:"event_type"`
}

type GroupEventStatsPayload struct {
	GroupID         string              `json:"group_id"`
	EventID         string              `json:"event_id"`
	GoingCount      int                 `json:"going_count"`
	InterestedCount int                 `json:"interested_count"`
	EventType       constants.EventType `json:"event_type"`
}
