package http

type MediaAlbumItem struct {
	Type     string          `json:"type"`
	Thumb    *MediaItemInfo  `json:"thumb,omitempty"`
	Detail   *MediaItemInfo  `json:"detail,omitempty"`
	Original *MediaItemInfo  `json:"original,omitempty"`
	Album    *MediaAlbumInfo `json:"album,omitempty"`
}

type MediaAlbumInfo struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}
