package strategy

import "github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/domain/IRepository/IStrategy"

func NewStrategyNotificationType(h1 *NotifCommentReply, h2 *NotifFriendAccept,
	h3 *NotifFriendRequest,
	h4 *NotifGroupInvite,
	h5 *NotifMention, h6 *NotiPostLike,
	h7 *NotiFSystemAlert) []IStrategy.IStrategyTypeNotificationType {
	return []IStrategy.IStrategyTypeNotificationType{
		h1, h2, h3, h4, h5, h6, h7,
	}
}
