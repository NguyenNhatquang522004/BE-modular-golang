package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/communityEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/res"
)

type StatsGroupEvent struct {
	events events.EventBus
}

func NewStatsGroupEvent(events events.EventBus) *StatsGroupEvent {
	return &StatsGroupEvent{
		events: events,
	}
}

func (s *StatsGroupEvent) Execute(ctx context.Context, req *req.StatsGroupEventRequest) (*res.FailedGroupEvent, error) {
	// Implement the logic for executing the StatsGroupEvent use case
	err := s.events.Publish(ctx, constants.TopicEventStats.String(), req.GroupID, req.EventType.String(), &communityEvent.GroupEventStatsPayload{
		GroupID:         req.GroupID,
		EventID:         req.EventID,
		GoingCount:      req.GoingCount,
		InterestedCount: req.InterestedCount,
		EventType:       req.EventType,
	})
	if err != nil {
		return &res.FailedGroupEvent{
			GroupID:      req.GroupID,
			EventID:      req.EventID,
			ErrorMessage: err,
		}, err
	}
	return nil, nil
}
