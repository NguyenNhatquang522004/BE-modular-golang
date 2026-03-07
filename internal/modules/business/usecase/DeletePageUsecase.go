package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/communicationEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/communityEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/contentEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/interactionEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/socialEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryMongodb"
)

type DeletePageUsecase struct {
	pageRepo       IRepositoryMongodb.IPagesRepository
	pageRole       IRepositoryMongodb.IPageRolesRepository
	pageFollow     IRepositoryMongodb.IPageFollowersRepository
	pageMetricRepo IRepositoryCassandra.IPageDailyMetricsRepository
	events         events.EventBus
}

func NewDeletePageUsecase(pageRepo IRepositoryMongodb.IPagesRepository, pageRole IRepositoryMongodb.IPageRolesRepository, pageFollow IRepositoryMongodb.IPageFollowersRepository, pageMetricRepo IRepositoryCassandra.IPageDailyMetricsRepository, events events.EventBus) *DeletePageUsecase {
	return &DeletePageUsecase{
		pageRepo:       pageRepo,
		pageRole:       pageRole,
		pageFollow:     pageFollow,
		pageMetricRepo: pageMetricRepo,
		events:         events,
	}
}

func (u *DeletePageUsecase) Execute(ctx context.Context, req *req.DeletePageRequest) (*res.FailedPageResponse, error) {
	datapage, err := u.pageRepo.GetPageByID(ctx, req.PageID)
	if err != nil {
		return &res.FailedPageResponse{
			PageID:       req.PageID,
			UserActionID: req.UserActionID,
			ErrorMessage: "Failed to retrieve page: " + err.Error(),
		}, err
	}
	if datapage == nil {
		return &res.FailedPageResponse{
			PageID:       req.PageID,
			UserActionID: req.UserActionID,
			ErrorMessage: "Page not found",
		}, nil
	}
	datauseraction, err := u.pageRole.GetPageRoleByPageIDAndUserID(ctx, datapage.ID.Hex(), req.UserActionID)
	if err != nil {
		return &res.FailedPageResponse{
			PageID:       req.PageID,
			UserActionID: req.UserActionID,
			ErrorMessage: "Failed to retrieve user role: " + err.Error(),
		}, err
	}
	if datauseraction == nil {
		return &res.FailedPageResponse{
			PageID:       req.PageID,
			UserActionID: req.UserActionID,
			ErrorMessage: "User has no role on this page",
		}, nil
	}
	if datapage.CreatorUserID != req.UserActionID {
		if datauseraction.Role != sharedEnums.RoleTypeAdmin {
			return &res.FailedPageResponse{
				PageID:       req.PageID,
				UserActionID: req.UserActionID,
				ErrorMessage: "User is not admin of this page",
			}, nil
		}
	}
	err = u.pageRepo.DeletePage(ctx, req.PageID)
	if err != nil {
		return &res.FailedPageResponse{
			PageID:       req.PageID,
			UserActionID: req.UserActionID,
			ErrorMessage: "Failed to delete page: " + err.Error(),
		}, err
	}
	err = u.pageRole.DeletePageRoleByPageID(ctx, req.PageID)
	if err != nil {
		return &res.FailedPageResponse{
			PageID:       req.PageID,
			UserActionID: req.UserActionID,
			ErrorMessage: "Failed to delete page roles: " + err.Error(),
		}, err
	}
	err = u.pageFollow.DeleteFollowerByPageID(ctx, req.PageID)
	if err != nil {
		return &res.FailedPageResponse{
			PageID:       req.PageID,
			UserActionID: req.UserActionID,
			ErrorMessage: "Failed to delete page followers: " + err.Error(),
		}, err
	}
	err = u.pageMetricRepo.DeletePageDailyMetricsByPageID(ctx, req.PageID)
	if err != nil {
		return &res.FailedPageResponse{
			PageID:       req.PageID,
			UserActionID: req.UserActionID,
			ErrorMessage: "Failed to delete page daily metrics: " + err.Error(),
		}, err
	}
	payload := &communicationEvent.DeletePrivateConversationGroupPayload{
		TargetID: req.PageID,
	}
	err = u.events.Publish(ctx, string(constants.TopicDeleteRelationTarget), req.PageID, constants.Deleted.String(), payload)
	if err != nil {
		return &res.FailedPageResponse{
			PageID:       req.PageID,
			UserActionID: req.UserActionID,
			ErrorMessage: "Failed to publish delete relation target event: " + err.Error(),
		}, err
	}
	payloadCommunity := &communityEvent.DeleteCommunityRelationTargetPayload{
		TargetID: req.PageID,
	}
	err = u.events.Publish(ctx, string(constants.TopicDeleteCommunityRelationTarget), req.PageID, constants.Deleted.String(), payloadCommunity)
	if err != nil {
		return &res.FailedPageResponse{
			PageID:       req.PageID,
			UserActionID: req.UserActionID,
			ErrorMessage: "Failed to publish delete community relation target event: " + err.Error(),
		}, err
	}
	payloadInteraction := &interactionEvent.DeleteInteractionRelationTargetPayload{
		TargetID: req.PageID,
	}
	err = u.events.Publish(ctx, string(constants.TopicDeleteInteractionRelationTarget), req.PageID, constants.Deleted.String(), payloadInteraction)
	if err != nil {
		return &res.FailedPageResponse{
			PageID:       req.PageID,
			UserActionID: req.UserActionID,
			ErrorMessage: "Failed to publish delete interaction relation target event: " + err.Error(),
		}, err
	}
	payloadContent := &contentEvent.DeleteContentRelationTargetPayload{
		TargetID: req.PageID,
	}
	err = u.events.Publish(ctx, string(constants.TopicDeleteContentRelationTarget), req.PageID, constants.Deleted.String(), payloadContent)
	if err != nil {
		return &res.FailedPageResponse{
			PageID:       req.PageID,
			UserActionID: req.UserActionID,
			ErrorMessage: "Failed to publish delete content relation target event: " + err.Error(),
		}, err
	}
	payloadSocial := &socialEvent.DeleteSocialRelationTargetPayload{
		TargetID: req.PageID,
	}
	err = u.events.Publish(ctx, string(constants.TopicDeleteSocialRelationTarget), req.PageID, constants.Deleted.String(), payloadSocial)
	if err != nil {
		return &res.FailedPageResponse{
			PageID:       req.PageID,
			UserActionID: req.UserActionID,
			ErrorMessage: "Failed to publish delete social relation target event: " + err.Error(),
		}, err
	}
	payloadMedia := &mediaEvent.DeleteMediaRelationTargetPayload{
		TargetID: req.PageID,
	}
	err = u.events.Publish(ctx, string(constants.TopicDeleteMediaRelationTarget), req.PageID, constants.Deleted.String(), payloadMedia)
	if err != nil {
		return &res.FailedPageResponse{
			PageID:       req.PageID,
			UserActionID: req.UserActionID,
			ErrorMessage: "Failed to publish delete media relation target event: " + err.Error(),
		}, err
	}
	return &res.FailedPageResponse{
		PageID:       req.PageID,
		UserActionID: req.UserActionID,
		ErrorMessage: "Delete page successfully",
	}, nil
}
