package sharedEnums

//go:generate enumer -type=NotificationType -json -transform=snake -trimprefix=Notif
type NotificationType int

const (
	NotifPostLike      NotificationType = iota // 'post_like'
	NotifCommentReply                          // 'comment_reply'
	NotifFriendRequest                         // 'friend_request'
	NotifFriendAccept                          // 'friend_accept'
	NotifGroupInvite                           // 'group_invite'
	NotifSystemAlert                           // 'system_alert'
	NotifMention                               // 'mention' (Được tag)
	NotifFollower                              // 'follow'
	NotifFriend                          // 'friend_request_accept'
)
