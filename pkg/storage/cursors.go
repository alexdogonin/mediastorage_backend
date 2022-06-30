package storage

import (
	"encoding/base64"
	"encoding/json"
)

type itemsCursor struct {
	TS    string `json:"uuid"`
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
	UUID     string `json:"uuid"`
	ItemUUID string `json:"item_uuid"`
	Limit    uint   `json:"limit"`
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
