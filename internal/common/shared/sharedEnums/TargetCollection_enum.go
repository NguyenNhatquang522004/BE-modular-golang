package sharedEnums

//go:generate enumer -type=TargetCollection -json -transform=snake -trimprefix=Collection
type TargetCollection int

const (
	CollectionPosts    TargetCollection = iota // 'posts'
	CollectionComments                         // 'comments'
)
