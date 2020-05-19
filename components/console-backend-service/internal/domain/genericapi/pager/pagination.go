package pager

import (
	"encoding/base64"
	"github.com/pkg/errors"
	"strconv"
)

type Page struct {
	StartCursor string
	EndCursor   string
	HasNextPage bool
}

func DecodeCursor(cursor string) (int, error) {
	if cursor == "" {
		return 0, nil
	}

	decodedValue, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0, errors.Wrap(err, "cursor is not correct")
	}

	offset, err := strconv.Atoi(string(decodedValue))
	if err != nil {
		return 0, errors.Wrap(err, "cursor is not correct")
	}
	if offset < 0 {
		return 0, errors.New("cursor is not correct")
	}

	return offset, nil
}

func EncodeNextCursor(offset, pageSize int) string {
	nextPage := pageSize + offset
	cursor := strconv.Itoa(nextPage)
	return base64.StdEncoding.EncodeToString([]byte(cursor))
}
