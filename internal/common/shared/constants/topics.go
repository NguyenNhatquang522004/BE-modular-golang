package constants

type EventType string

const (
	Updated EventType = "UPDATED"
	Deleted EventType = "DELETED"
	Created EventType = "CREATED"
	None    EventType = "NONE"
	Failed  EventType = "FAILED"
)

type ProcessStatus string

const (
	StatusUndefined  ProcessStatus = "UNDEFINED"
	StatusProcessing ProcessStatus = "PROCESSING"
	StatusSuccess    ProcessStatus = "SUCCESS"
	StatusFailed     ProcessStatus = "FAILED"
	StatusSkipped    ProcessStatus = "SKIPPED" // Dùng khi message bị trùng (Idempotent) hoặc không hợp lệ để xử lý

)

type TopicName string

const (
	//DQL chung
	TopicDLQ TopicName = "dlq.events"
	//identity
	TopicUserSettings TopicName = "identity.user_settings.events"
	//socials
	TopicFriendship                 TopicName = "social.friendship.events"
	TopicBlockUser                  TopicName = "social.block_user.events"
	TopicFollowUser                 TopicName = "social.follow_user.events"
	TopicProfile                    TopicName = "social.profile.events"
	TopicDeleteSocialRelationTarget TopicName = "social.delete_relation_target.events" // chưa làm gì cả

	//content
	TopicPost                        TopicName = "content.post.events"
	TopicSharePost                   TopicName = "content.post.share.events"
	TopicDeleteContentRelationTarget TopicName = "content.delete_relation_target.events" // chưa làm gì cả
	TopicPostStats                   TopicName = "content.post_stats.events"

	// notification
	TopicUserNotificationSettings TopicName = "notification.create_user_notification_settings.events"
	TopicSendNotificationType     TopicName = "notification.send_notification.events"

	//interaction
	TopicEdit                            TopicName = "interaction.edit_comment.events"
	TopicUserBookmark                    TopicName = "interaction.user_bookmark.events"
	TopicCommentStats                    TopicName = "interaction.comment_stats.events"
	TopicComment                         TopicName = "interaction.comment_reaction.events"
	TopicEntityReaction                  TopicName = "interaction.entity_reaction.events"
	TopicDeleteInteractionRelationTarget TopicName = "interaction.delete_relation_target.events" // chưa làm gì cả

	// media
	TopicMediaAsset                TopicName = "media.asset.events"
	TopicMeiaAssetHandleMetadata   TopicName = "media.asset_handle_metadata.events"
	TopicAlbumStats                TopicName = "media.album_stats.events"
	TopicStoryStats                TopicName = "media.story_stats.events"
	TopicReelStats                 TopicName = "media.reel_stats.events"
	TopicLiveSessionStats          TopicName = "media.live_session.events"
	TopicStartStopLive             TopicName = "media.start_stop_live.events" //refactor sau
	TopicCommentLive               TopicName = "media.comment_live.events"
	TopicDeleteMediaRelationTarget TopicName = "media.delete_relation_target.events" // chưa làm gì cả
	TopicStory                     TopicName = "media.story.events"
	TopicReel                      TopicName = "media.reel.events"
	TopicAlbum                     TopicName = "media.album.events"
	TopicMusic                     TopicName = "media.music.events"
	TopicMusicStats                TopicName = "media.music_stats.events"
	TopicArtist                    TopicName = "media.artist.events"
	TopicArtistStats               TopicName = "media.artist_stats.events"
	TopicLiveSession               TopicName = "media.live_session.events"
	// communication
	TopicMessage                 TopicName = "communication.message.events"
	TopicStatsMessage            TopicName = "communication.stats_message.events"
	TopicDeleteRelationTarget    TopicName = "business.delete_relation_target.events"
	TopicConversation            TopicName = "communication.conversation.events"
	TopicStatsConversation       TopicName = "communication.stats_conversation.events"
	TopicConversationParticipant TopicName = "communication.participant.events"
	TopicCallLog                 TopicName = "communication.call_log.events"
	TopicRelpyStory              TopicName = "communication.reply_story.events"
	//community
	TopicGroupStats                    TopicName = "community.group_stats.events"
	TopicGroup                         TopicName = "community.group.events"
	TopicGroupMember                   TopicName = "community.group_member.events"
	TopicGroupEvent                    TopicName = "community.group_event.events"
	TopicGroupQA                       TopicName = "community.group_qa.events"
	TopicGroupFile                     TopicName = "community.group_file.events"
	TopicDownloadGroupFile             TopicName = "community.download_group_file.events"
	TopicCommunityStats                TopicName = "community.community_stats.events"
	TopicDeleteCommunityRelationTarget TopicName = "business.delete_relation_target.events" // chưa làm gì cả, dùng chung với communication

	// business
	TopicPage             TopicName = "business.page.events"
	TopicPageRole         TopicName = "business.page_role.events"
	TopicFollowerPage     TopicName = "business.follower_page.events"
	TopicStatsPage        TopicName = "business.stats_page.events"
	TopicDailyMetricsPage TopicName = "business.daily_metrics_page.events"
)

type TopicConfig struct {
	Name       TopicName
	Partitions int
}

func (t ProcessStatus) String() string {
	return string(t)
}
func (t TopicName) String() string {
	return string(t)
}

func (t EventType) String() string {
	return string(t)
}

var SocialTopics = []TopicConfig{
	{Name: TopicFriendship, Partitions: 6}, // Gom 3 cái Create/Update/Delete vào 1
	{Name: TopicBlockUser, Partitions: 6},
	{Name: TopicFollowUser, Partitions: 3},
	{Name: TopicUserNotificationSettings, Partitions: 3},
	{Name: TopicSendNotificationType, Partitions: 3},
	{Name: TopicDeleteRelationTarget, Partitions: 3},
	{Name: TopicDailyMetricsPage, Partitions: 3},
}
