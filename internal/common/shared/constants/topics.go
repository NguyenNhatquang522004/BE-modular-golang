package constants

type EventType string

const (
	Updated EventType = "UPDATED"
	Deleted EventType = "DELETED"
	Created EventType = "CREATED"
	None    EventType = "NONE"
	Failed  EventType = "FAILED"
)

type TopicName string

const (
	//socials
	TopicFriendship                 TopicName = "social.friendship.events"
	TopicBlock                      TopicName = "social.block.events"
	TopicFollow                     TopicName = "social.follow.events"
	TopicDeleteSocialRelationTarget TopicName = "social.delete_relation_target.events" // chưa làm gì cả

	//content
	TopicContentPostPublish              TopicName = "content.post.publish.events"
	TopicContentPostPublishMediaAssets   TopicName = "content.post.publish_media_assets.events"
	ContentPostPublishNotificationFriend TopicName = "content.post.publish_notification_friend.events"
	ContentPostPublishNotificationTag    TopicName = "content.post.publish_notification_tag.events"
	TopicSharePost                       TopicName = "content.post.share.events"
	TopicDeleteContentRelationTarget     TopicName = "content.delete_relation_target.events" // chưa làm gì cả

	// notification
	TopicCreateUserNotificationSettings TopicName = "notification.create_user_notification_settings.events"
	TopicSendNotificationType           TopicName = "notification.send_notification.events"

	//interaction
	TopicReactComment                    TopicName = "interaction.react_comment.events"
	TopicReactPost                       TopicName = "interaction.react_post.events"
	TopicCommentPost                     TopicName = "interaction.comment_post.events"
	TopicCounterPost                     TopicName = "interaction.counter_post.events"
	TopicCounterComment                  TopicName = "interaction.counter_comment.events"
	TopicEntityReaction                  TopicName = "interaction.entity_reaction.events"
	TopicDeleteInteractionRelationTarget TopicName = "interaction.delete_relation_target.events" // chưa làm gì cả

	// media
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
	TopicStatsPage TopicName = "business.stats_page.events"
)

type TopicConfig struct {
	Name       TopicName
	Partitions int
}

func (t TopicName) String() string {
	return string(t)
}

func (t EventType) String() string {
	return string(t)
}

var SocialTopics = []TopicConfig{
	{Name: TopicFriendship, Partitions: 6}, // Gom 3 cái Create/Update/Delete vào 1
	{Name: TopicBlock, Partitions: 6},
	{Name: TopicFollow, Partitions: 3},
	{Name: TopicContentPostPublish, Partitions: 3},
	{Name: TopicContentPostPublishMediaAssets, Partitions: 3},
	{Name: TopicCreateUserNotificationSettings, Partitions: 3},
	{Name: TopicSendNotificationType, Partitions: 3},
	{Name: TopicDeleteRelationTarget, Partitions: 3},
}
