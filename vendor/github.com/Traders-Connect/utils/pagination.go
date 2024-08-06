package utils

import (
	b64 "encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Cursor is a structure maintaing information about the cursor
type Cursor struct {
	ID        uint
	Name      string
	TimeStamp time.Time
}

// ToBase64String converts a cursor to a base64 encoded string
func (c *Cursor) ToBase64String() string {
	cursor := fmt.Sprintf("%d/%s/%s", c.ID, c.Name, c.TimeStamp)
	return b64.StdEncoding.EncodeToString([]byte(cursor))
}

// ParseCursor returns a cursor structure from a base64 encoded string
func ParseCursor(cursor string) (*Cursor, error) {
	fromID, err := b64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, fmt.Errorf("error decoding pagination cursor")
	}

	fields := strings.Split(string(fromID), "/")
	if len(fields) < 3 {
		return nil, fmt.Errorf("invalid pagination cursor")
	}
	cID, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("invalid ID")
	}

	t, err := time.Parse(UTCLayout, fields[len(fields)-1])
	if err != nil {
		return nil, fmt.Errorf("can't parse timestamp")
	}

	return &Cursor{ID: uint(cID), Name: strings.Join(fields[1:len(fields)-1], "/"), TimeStamp: t}, nil
}
