package utils

import (
	b64 "encoding/base64"
	"fmt"
	"strings"
	"time"
)

// UUIDCursor is a structure maintaing information about the uuid cursor
type UUIDCursor struct {
	ID        string
	TimeStamp time.Time
}

// ToBase64String converts a cursor to a base64 encoded string
func (u *UUIDCursor) ToBase64String() string {
	cursor := fmt.Sprintf("%s,%s", u.ID, u.TimeStamp.Format(time.RFC3339Nano))
	return b64.StdEncoding.EncodeToString([]byte(cursor))
}

// ParseUUIDCursor returns a cursor structure from a base64 encoded string
func ParseUUIDCursor(cursor string) (*UUIDCursor, error) {
	fromID, err := b64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, fmt.Errorf("error decoding pagination cursor")
	}

	fields := strings.Split(string(fromID), ",")
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid pagination cursor")
	}

	uuid := fields[0]

	t, err := time.Parse(time.RFC3339Nano, fields[len(fields)-1])
	if err != nil {
		return nil, fmt.Errorf("can't parse timestamp")
	}

	return &UUIDCursor{ID: uuid, TimeStamp: t}, nil
}
