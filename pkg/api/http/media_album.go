package http

type MediaAlbum struct {
	Name  string
	Items []MediaAlbumItem `json:"items"`
}

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

type MediaAlbumItemType uint