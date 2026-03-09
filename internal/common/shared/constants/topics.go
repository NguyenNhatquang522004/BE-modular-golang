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
	TopicCommentStats                    TopicName = "interaction.comment_stats.events"
	TopicEntityReaction                  TopicName = "interaction.entity_reaction.events"
	TopicDeleteInteractionRelationTarget TopicName = "interaction.delete_relation_target.events" // chưa làm gì cả

	// media
	TopicMediaAsset     TopicName = "media.asset.events"
	TopicReactAlbum     TopicName = "media.react_album.events"
	TopicViewCountStory TopicName = "media.view_count_story.events"
	TopicReactStory     TopicName = "media.react_story.events"
	TopicReplyStory     TopicName = "media.rely_story.events"
	TopicReactReel      TopicName = "media.react_reel.events"
	TopicCounterReel    TopicName = "media.counter_reel.events"
	TopicReactLive      TopicName = "media.react_live.events"
	TopicStartStopLive  TopicName = "media.start_stop_live.events"
	TopicCommentLive    TopicName = "media.comment_live.events"
	TopicCounterLive    TopicName = "media.counter_live.events"
	TopicDelete
	TopicDeleteMediaRelationTarget TopicName = "media.delete_relation_target.events" // chưa làm gì cả
	// communication
	TopicMessage              TopicName = "communication.message.events"
	TopicStateMessage         TopicName = "communication.state_message.events"
	TopicReactMessage         TopicName = "communication.react_message.events"
	TopicDeleteRelationTarget TopicName = "business.delete_relation_target.events"
	//community
	TopicGroupStats                    TopicName = "community.group_stats.events"
	TopicEventStats                    TopicName = "community.event_stats.events"
	TopicDownloadGroupFile             TopicName = "community.download_group_file.events"
	TopicDeleteCommunityRelationTarget TopicName = "business.delete_relation_target.events" // chưa làm gì cả, dùng chung với communication

	// business
	TopicStatsPage        TopicName = "business.stats_page.events"
	TopicFollowerPage     TopicName = "business.follower_page.events"
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
