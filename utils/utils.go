package utils

import (
	"fmt"
	"strconv"
	"time"
)

// StringSet is a set(data structure) for string
type StringSet map[string]bool

// Contains check if element is in the set
func (set StringSet) Contains(e string) bool {
	val, ok := set[e]
	return ok && val
}

// Add adds element into the set
func (set StringSet) Add(e string) {
	set[e] = true
}

// Remove removes element in the set
func (set StringSet) Remove(e string) {
	delete(set, e)
}

// Converts string-formatted timestamp to UNIX Int64
func TimestampAsInt64(t string) (int64, error) {
	timestamp, err := time.Parse(time.RFC3339Nano, t)
	if err != nil {
		return 0, err
	}
	return timestamp.Unix(), err
}

// ExtractTimestamp parses ISO8601 or RFC3339Nano formatted time string to int64 UNIX timestamp.
func ExtractTimestamp(obj map[string]interface{}, key string) (int64, bool) {
	timestring, ok := obj[key].(string)
	if !ok {
		return 0, false
	}
	timestamp, err := time.Parse(time.RFC3339Nano, timestring)
	if err != nil {
		fmt.Println(err)
		return 0, false
	}
	return timestamp.Unix(), true
}

// ExtractInt64 parses a value in the given map to uint64.
func ExtractInt64(obj map[string]interface{}, key string) (int64, bool) {
	num, ok := obj[key].(int64)
	if !ok {
		nstring, ok := obj[key].(string)
		if !ok {
			return 0, false
		}
		n, err := strconv.ParseInt(nstring, 10, 64)
		if err != nil {
			return 0, false
		}
		num = n
	}
	return num, true
}

// ExtractUInt64 parses a value in the given map to uint64.
func ExtractUInt64(obj map[string]interface{}, key string) (uint64, bool) {
	num, ok := obj[key].(uint64)
	if !ok {
		nstring, ok := obj[key].(string)
		if !ok {
			return 0, false
		}
		n, err := strconv.ParseUint(nstring, 10, 64)
		if err != nil {
			return 0, false
		}
		num = n
	}
	return num, true
}

// ExtractInt32 parses a value in the given map to int32.
func ExtractInt32(obj map[string]interface{}, key string) (int, bool) {
	num, ok := obj[key].(int)
	if !ok {
		nstring, ok := obj[key].(string)
		if !ok {
			return 0, false
		}
		n, err := strconv.ParseInt(nstring, 10, 32)
		if err != nil {
			return 0, false
		}
		num = int(n)
	}
	return num, true
}

// ExtractFloat64 parses a value in the given map to float64.
func ExtractFloat64(obj map[string]interface{}, key string) (float64, bool) {
	num, ok := obj[key].(float64)
	if !ok {
		nstring, ok := obj[key].(string)
		if !ok {
			return 0, false
		}
		n, err := strconv.ParseFloat(nstring, 64)
		if err != nil {
			return 0, false
		}
		num = float64(n)
	}
	return num, true
}

// ExtractFloat32 parses a value in the given map to float32.
func ExtractFloat32(obj map[string]interface{}, key string) (float32, bool) {
	num, ok := obj[key].(float32)
	if !ok {
		nstring, ok := obj[key].(string)
		if !ok {
			return 0, false
		}
		n, err := strconv.ParseFloat(nstring, 32)
		if err != nil {
			return 0, false
		}
		num = float32(n)
	}
	return num, true
}
