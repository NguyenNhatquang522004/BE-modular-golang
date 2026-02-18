package blockStrategy

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IStrategy"
)

// func NewProviderBlockFull(blockRepo IRepositoryPostgres.IBlockRepository,
// 	blockProducer graph.IBlockMessage,
// 	friendshipRepo IRepositoryPostgres.IFriendshipsRepository, friendshipProducer graph.IFriendshipMessage,
// 	followRepo IRepositoryPostgres.IFollowersRepository,
// 	followProducer graph.IFollowMessage) *BlockFull {
// 	return &BlockFull{
// 		blockRepo:          blockRepo,
// 		blockProducer:      blockProducer,
// 		friendshipRepo:     friendshipRepo,
// 		friendShipProducer: friendshipProducer,
// 		followRepo:         followRepo,
// 		followProducer:     followProducer,
// 	}
// }

// func NewProviderBlockProfile(blockRepo IRepositoryPostgres.IBlockRepository,
// 	blockProducer graph.IBlockMessage,
// 	friendshipRepo IRepositoryPostgres.IFriendshipsRepository, friendshipProducer graph.IFriendshipMessage,
// 	followRepo IRepositoryPostgres.IFollowersRepository,
// 	followProducer graph.IFollowMessage) *BlockProfile {
// 	return &BlockProfile{
// 		blockRepo:          blockRepo,
// 		blockProducer:      blockProducer,
// 		friendshipRepo:     friendshipRepo,
// 		friendShipProducer: friendshipProducer,
// 		followRepo:         followRepo,
// 		followProducer:     followProducer,
// 	}
// }

// func NewProviderBlockChat(blockRepo IRepositoryPostgres.IBlockRepository,
// 	blockProducer graph.IBlockMessage,
// 	friendshipRepo IRepositoryPostgres.IFriendshipsRepository, friendshipProducer graph.IFriendshipMessage,
// 	followRepo IRepositoryPostgres.IFollowersRepository,
// 	followProducer graph.IFollowMessage) *BlockChat {
// 	return &BlockChat{
// 		blockRepo:          blockRepo,
// 		blockProducer:      blockProducer,
// 		friendshipRepo:     friendshipRepo,
// 		friendShipProducer: friendshipProducer,
// 		followRepo:         followRepo,
// 		followProducer:     followProducer,
// 	}
// }

// func NewProviderBlockNone(blockRepo IRepositoryPostgres.IBlockRepository,
// 	blockProducer graph.IBlockMessage,
// 	friendshipRepo IRepositoryPostgres.IFriendshipsRepository, friendshipProducer graph.IFriendshipMessage,
// 	followRepo IRepositoryPostgres.IFollowersRepository,
// 	followProducer graph.IFollowMessage) *BlockNone {
// 	return &BlockNone{
// 		blockRepo:          blockRepo,
// 		blockProducer:      blockProducer,
// 		friendshipRepo:     friendshipRepo,
// 		friendShipProducer: friendshipProducer,
// 		followRepo:         followRepo,
// 		followProducer:     followProducer,
// 	}
// }
func NewProviderBlockStrategy(h1 *BlockFull, h2 *BlockProfile, h3 *BlockChat, h4 *BlockNone) []IStrategy.IBlockStrategy {
	return []IStrategy.IBlockStrategy{h1, h2, h3, h4}
}
