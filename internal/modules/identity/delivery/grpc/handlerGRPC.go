package grpc

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/pkg/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IdentityHandlerGRPC struct {
	pb.UnimplementedIdentityServiceServer
	userSettingRepo IRepositoryMongodb.IUserSettingRepository
}

func NewIdentityHandlerGRPC(userSettingRepo IRepositoryMongodb.IUserSettingRepository) *IdentityHandlerGRPC {
	return &IdentityHandlerGRPC{
		userSettingRepo: userSettingRepo,
	}
}

func (h *IdentityHandlerGRPC) GetUserOTP(ctx context.Context, req *emptypb.Empty) (*pb.UserIDResponse, error) {
	// 1. Bóc Metadata từ Context (Incoming)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "Thiếu metadata")
	}

	// 2. Lấy giá trị động ra bằng cái Key mà A đã nhét vào
	userIDs := md["x-user-id"]
	if len(userIDs) == 0 {
		return nil, status.Error(codes.Unauthenticated, "Không tìm thấy x-user-id trong context")
	}

	// Đây chính là dữ liệu "động" mà A truyền sang!
	dynamicID := userIDs[0]

	// 3. Bây giờ B có ID động rồi, B có thể gọi Database để lấy data thật
	// user, err := h.usecase.FindUserInDB(dynamicID)

	// Trả về cho A
	return &pb.UserIDResponse{
		UserId: dynamicID, // Trả lại đúng cái ID động đó
	}, nil
}
func (h *IdentityHandlerGRPC) GetUserSettingByID(ctx context.Context, req *pb.UserSettingIDRequest) (*pb.UserSettingResponse, error) {
	datauser, err := h.userSettingRepo.GetUserSettings(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Lỗi khi lấy cài đặt người dùng")
	}
	if datauser == nil {
		return nil, status.Error(codes.Internal, "Lỗi khi lấy cài đặt người dùng")
	}
	// Chuyển đổi datauser sang UserSettingResponse nếu cần thiết
	response := &pb.UserSettingResponse{
		// Điền các trường dữ liệu từ datauser vào đây
		Id:                         datauser.ID.Hex(),
		UserId:                     datauser.User_ID,
		DefaultPostAudience:        datauser.Default_Post_Audience.String(),
		DefaultStoryAudience:       datauser.Default_Story_Audience.String(),
		AllowFriendRequestFrom:     datauser.Allow_Friend_Request_From.String(),
		AllowFriendListViewFrom:    datauser.Allow_Friend_List_View_From.String(),
		AllowSearchEngineIndexing:  datauser.Allow_Search_Engine_Indexing,
		Allow_Profile_View_From:    datauser.Allow_Profile_View_From,
		AllowTimelinePostingFrom:   datauser.Allow_Timeline_Posting_From.String(),
		ReviewTagsEnabled:          datauser.Review_Tags_Enabled,
		ReviewTimelinePostsEnabled: datauser.Review_Timeline_Posts_Enabled,
		Notifications: &pb.NotificationSettings{
			EmailFrequency:   datauser.Notifications.EmailFrequency.String(),
			PushInteractions: datauser.Notifications.PushInteractions,
			PushFriends:      datauser.Notifications.PushFriends,
			PushGroups:       datauser.Notifications.PushGroups,
			PushEvents:       datauser.Notifications.PushEvents,
			PushBirthdays:    datauser.Notifications.PushBirthdays,
		},
		UpdatedAt: timestamppb.New(datauser.Updated_At), // Định dạng thời gian theo ISO 8601
	}
	return response, nil
}
