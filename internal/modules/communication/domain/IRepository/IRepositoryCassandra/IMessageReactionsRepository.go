package IRepositoryCassandra

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/cassandraErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/entity"
)

// IMessageReactionsRepository định nghĩa hợp đồng truy cập dữ liệu cho bảng message_reactions.
//
// Schema Cassandra:
//
//	PRIMARY KEY ((conversation_id, message_id), user_id)
//
// Nguyên tắc thiết kế:
//   - Mọi query ĐỀU PHẢI chỉ định đủ Composite Partition Key (conversation_id + message_id)
//     để tránh Full Cluster Scan.
//   - Xóa toàn bộ reaction của 1 tin nhắn (DeleteReactionsByMessageID) là Partition Delete
//     — thao tác rất nhanh trong Cassandra, không cần Worker Pool.
//   - Bulk operations dùng Worker Pool để thực thi song song, trả về Partial-Success.
type IMessageReactionsRepository interface {

	// =========================================================================
	// CREATE
	// =========================================================================

	// CreateReaction tạo mới 1 reaction. Upsert theo cơ chế của Cassandra.
	// Nếu user đã react trước đó, reaction_code sẽ được ghi đè (last-write-wins).
	CreateReaction(ctx context.Context, reaction *entity.MessageReaction) error

	// CreateBulkReactions tạo nhiều reaction song song qua Worker Pool.
	// Trả về (số thành công, danh sách lỗi chi tiết, lỗi tổng hợp).
	CreateBulkReactions(ctx context.Context, reactions []*entity.MessageReaction) (int64, []*cassandraErrors.MessageReactionBulkError, error)

	// =========================================================================
	// UPDATE
	// =========================================================================

	// UpdateReaction cập nhật reaction_code của 1 reaction đã tồn tại.
	// Yêu cầu đủ Primary Key: (conversation_id, message_id, user_id).
	UpdateReaction(ctx context.Context, reaction *entity.MessageReaction) error

	// UpdateBulkReactions cập nhật nhiều reaction song song qua Worker Pool.
	// Trả về (số thành công, danh sách lỗi chi tiết, lỗi tổng hợp).
	UpdateBulkReactions(ctx context.Context, reactions []*entity.MessageReaction) (int64, []*cassandraErrors.MessageReactionBulkError, error)

	// =========================================================================
	// DELETE
	// =========================================================================
	DeleteReactionByConversationID(ctx context.Context, conversationID string) error
	// DeleteReaction xóa chính xác 1 reaction theo đủ Primary Key.
	// (conversation_id, message_id) → Partition | user_id → Clustering Key.
	DeleteReaction(ctx context.Context, conversationID string, messageID string, userID string) error

	// DeleteBulkReactions xóa reaction của nhiều user trên cùng 1 tin nhắn song song.
	// Trả về (số thành công, danh sách lỗi chi tiết, lỗi tổng hợp).
	DeleteBulkReactions(ctx context.Context, conversationID string, messageID string, userIDs []string) (int64, []*cassandraErrors.MessageReactionBulkError, error)

	// DeleteReactionsByMessageID xóa TOÀN BỘ reaction của 1 tin nhắn (Partition Delete).
	// Đây là thao tác cực nhanh trong Cassandra — chỉ cần ghi 1 Tombstone-Partition.
	// Dùng khi thu hồi (revoke) hoặc xóa tin nhắn gốc.
	DeleteReactionsByMessageID(ctx context.Context, conversationID string, messageID string) error

	// DeleteBulkReactionsByMessageID xóa toàn bộ reaction của nhiều tin nhắn song song.
	// Mỗi messageID sẽ kích hoạt 1 Partition Delete độc lập trên Worker Pool.
	// Trả về (số partition xóa thành công, danh sách lỗi chi tiết, lỗi tổng hợp).
	DeleteBulkReactionsByMessageID(ctx context.Context, conversationID string, messageIDs []string) (int64, []*cassandraErrors.MessageReactionBulkError, error)

	// =========================================================================
	// READ
	// =========================================================================

	// GetReactionsByMessageID lấy danh sách reaction của 1 tin nhắn với phân trang cursor.
	// Chỉ cần Partition Key đầy đủ (conversation_id + message_id) → truy vấn nhanh.
	GetReactionsByMessageID(ctx context.Context, conversationID string, messageID string, cursor string, limit int) (*dto.PaginationRes, error)

	// GetReactionByUser lấy đúng 1 reaction của 1 user trên 1 tin nhắn bằng đủ Primary Key.
	// Trả về nil, nil nếu user chưa thả reaction.
	GetReactionByUser(ctx context.Context, conversationID string, messageID string, userID string) (*entity.MessageReaction, error)
}
