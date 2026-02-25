package constants

type EventType string

const (
	Updated EventType = "UPDATED"
	Deleted EventType = "DELETED"
	Created EventType = "CREATED"
	None    EventType = "NONE"
)

type TopicName string

const (
	//socials
	TopicFriendship TopicName = "social.friendship.events"
	TopicBlock      TopicName = "social.block.events"
	TopicFollow     TopicName = "social.follow.events"

	//content
	TopicContentPostPublish              TopicName = "content.post.publish.events"
	TopicContentPostPublishMediaAssets   TopicName = "content.post.publish_media_assets.events"
	ContentPostPublishNotificationFriend TopicName = "content.post.publish_notification_friend.events"
	ContentPostPublishNotificationTag    TopicName = "content.post.publish_notification_tag.events"

	// notification
	TopicCreateUserNotificationSettings TopicName = "notification.create_user_notification_settings.events"
	TopicSendNotificationType           TopicName = "notification.send_notification.events"
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
	// Thêm các topic khác vào đây
}
