package utils

import (
	"strconv"
	"net/http"
	"gopkg.in/guregu/null.v3"
)

const (
	DefaultCursor = -1
	//200
	DefaultPageSize = 5
	MaxPageSize     = 1000
)

type Page struct {
	Data          []interface{} `json:"data"`
	NextCursor    null.Int      `json:"next_cursor"`
	NextCursorStr string        `json:"next_cursor_str"`
	PrevCursor    int64         `json:"previous_cursor"`
	PrevCursorStr string        `json:"previous_cursor_str"`
}

func EmptyPage() Page {
	data := make([]interface{}, 0)
	return Page{
		Data:          data,
		PrevCursor:    DefaultCursor,
		PrevCursorStr: strconv.FormatInt(DefaultCursor, 10),
		NextCursor:    null.IntFrom(DefaultCursor),
		NextCursorStr: strconv.FormatInt(DefaultCursor, 10)}
}

func MakePage(count int, data []interface{}, cursor int64, lastId int64) Page {
	var nextCursor null.Int
	var nextCursorStr string
	if len(data) < count {
		nextCursor = null.NewInt(0, false)
	} else {
		nextCursor = null.IntFrom(lastId)
		nextCursorStr = strconv.FormatInt(lastId, 10)
	}

	return Page{
		Data:          data,
		PrevCursor:    cursor,
		PrevCursorStr: strconv.FormatInt(cursor, 10),
		NextCursor:    nextCursor,
		NextCursorStr: nextCursorStr}
}

func GetCursor(r *http.Request) int64 {
	cursorString := r.FormValue("cursor")
	if cursorString == "" {
		return DefaultCursor
	}
	if !IsParseableAsInt64(cursorString) {
		return DefaultCursor
	}
	cursor := StringToInt64(cursorString)
	return cursor
}

func GetSort(r *http.Request) string {
	// TODO actually parse it
	return "ASC"
}

func GetCount(r *http.Request) int {
	countString := r.FormValue("count")
	if countString == "" {
		return DefaultPageSize
	}
	count, _ := strconv.Atoi(countString)
	if count > MaxPageSize {
		count = MaxPageSize
	}
	if count < 1 {
		count = 1
	}
	return count
}