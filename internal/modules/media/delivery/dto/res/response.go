package res

type CreateAlbumResponse struct {
	Fail               []*FailedMediaAssetsResponse `json:"fail"`
	AlbumID            string                       `json:"album_id"`
	AlbumsErrorMessage string                       `json:"albums_error_message"`
}

type FailedMediaAssetsResponse struct {
	MediaAssetsId string `json:"media_asset_id"`
	ErrorMessage  string `json:"error_message"`
}

type UpdateAlbumResponse struct {
	AlbumID            string `json:"album_id"`
	AlbumsErrorMessage string `json:"albums_error_message"`
}

type DeleteAlbumResponse struct {
	AlbumID            string                       `json:"album_id"`
	AlbumsErrorMessage string                       `json:"albums_error_message"`
	Fail               []*FailedMediaAssetsResponse `json:"fail"`
}

type ReactAlbumResponse struct {
	AlbumID            string `json:"album_id"`
	AlbumsErrorMessage string `json:"albums_error_message"`
	EventType          string `json:"event_type"`
}

type FailedStoryResponse struct {
	StoryID      string `json:"story_id"`
	UserID       string `json:"user_id"`
	ErrorMessage string `json:"error_message"`
}
