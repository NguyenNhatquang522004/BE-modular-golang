package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/businessEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/res"
)

type StatsPageUsecase struct {
	events events.EventBus
}

func NewStatsPageUsecase(events events.EventBus) *StatsPageUsecase {
	return &StatsPageUsecase{
		events: events,
	}
}

func (u *StatsPageUsecase) Execute(ctx context.Context, req *req.StatsPageRequest) (*res.FailedPageResponse, error) {
	// Thực hiện logic xử lý thống kê trang ở đây, ví dụ: gọi repository để cập nhật thống kê

	// Sau khi xử lý xong, bạn có thể publish một sự kiện nếu cần thiết
	payload := &businessEvent.StatsPagePayload{
		PageID:         req.PageID,
		UserID:         req.UserID,
		FollowersCount: req.FollowersCount,
		LikesCount:     req.LikesCount,
		RatingScore:    req.RatingScore,
		ReviewCount:    req.ReviewCount,
		CreatedAt:      req.CreatedAt,
	}
	err := u.events.Publish(ctx, constants.TopicStatsPage.String(), req.PageID, constants.Created.String(), payload)
	if err != nil {
		return &res.FailedPageResponse{
			PageID:       req.PageID,
			UserActionID: req.UserID,
			ErrorMessage: "Failed to publish stats page event: " + err.Error(),
		}, err
	}
	// Trả về phản hồi thành công hoặc lỗi tùy thuộc vào kết quả xử lý
	return &res.FailedPageResponse{
		PageID:       req.PageID,
		UserActionID: req.UserID,
		ErrorMessage: "Failed to update stats page",
	}, nil
}
