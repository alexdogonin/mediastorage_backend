package cache

import (
	"encoding/base64"
	"encoding/json"
)

type itemsCursor struct {
	UUID  string `json:"uuid"`
	Limit uint   `json:"limit"`
}

func (c *itemsCursor) Parse(cursor string) error {
	data, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, c)
}

func (c itemsCursor) String() string {
	data, _ := json.Marshal(c)

	return base64.RawURLEncoding.EncodeToString(data)
}

type albumsCursor struct {
	AlbumUUID string `json:"album_uuid"`
	Offset    uint   `json:"offset"`
	Limit     uint   `json:"limit"`
}

func (c *albumsCursor) Parse(cursor string) error {
	data, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, c)
}

func (c albumsCursor) String() string {
	data, _ := json.Marshal(c)

	return base64.RawURLEncoding.EncodeToString(data)
}
