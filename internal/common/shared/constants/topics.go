package constants

type EventType string

const (
	Updated EventType = "FRIENDSHIP_UPDATED"
	Deleted EventType = "FRIENDSHIP_DELETED"
)

type TopicName string

const (
	TopicFriendship TopicName = "social.friendship.events"
	TopicBlock      TopicName = "social.block.events"
	TopicFollow     TopicName = "social.follow.events"
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
	// Thêm các topic khác vào đây
}
