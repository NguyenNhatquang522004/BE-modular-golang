package friendshipstrategy

import "github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IStrategy"

func NewProviderfriendshipAccepted(h1 *FriendShipAccepted,
	h2 *FriendshipDeclined,
	h3 *FriendshipPending,
	h4 *FriendshipBlocked) []IStrategy.IFriendshipStrategy {
	return []IStrategy.IFriendshipStrategy{h1, h2, h3, h4}
}
