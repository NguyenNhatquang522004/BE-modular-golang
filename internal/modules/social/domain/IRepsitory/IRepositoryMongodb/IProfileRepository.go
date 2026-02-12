package IRepositoryMongodb


type IProfileRepositoryMongodb interface {
	// Define methods for profile repository
	CreateProfile(profileData interface{}) error
	GetProfileByID(profileID string) (interface{}, error)
	UpdateProfile(profileID string, updateData interface{}) error
	DeleteProfile(profileID string) error
}	