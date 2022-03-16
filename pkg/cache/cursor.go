package cache

import (
	"encoding/base64"
	"encoding/json"
)

type cursor struct {
	UUID  string `json:"uuid"`
	Limit uint   `json:"limit"`
}

func (c *cursor) Parse(cursor string) error {
	data, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, c)
}

func (c cursor) String() string {
	data, _ := json.Marshal(c)

	return base64.RawURLEncoding.EncodeToString(data)
}
