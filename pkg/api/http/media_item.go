package http

type MediaItem struct {
	UUID        string         `json:"uuid"`
	ThumbURL    string         `json:"thumb_url,omitempty"`
	DetailURL   string         `json:"detail_url,omitempty"`
	OriginalURL string         `json:"original_url,omitempty"`
	Thumb       *MediaItemInfo `json:"thumb,omitempty"`
	Detail      *MediaItemInfo `json:"detail,omitempty"`
	Original    *MediaItemInfo `json:"original,omitempty"`
}

type MediaItemInfo struct {
	URL    string `json:"url"`
	Width  uint   `json:"width"`
	Height uint   `json:"height"`
}
