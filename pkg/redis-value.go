package main

import (
	"errors"
	"strconv"
)

// errUnexpectedHotkeysReply is returned when HOTKEYS GET does not reply with a map.
var errUnexpectedHotkeysReply = errors.New("unexpected HOTKEYS GET reply")

/**
 * redisValueToString renders a RESP value as text.
 *
 * radix decodes bulk strings as []byte and integers as int64, and RESP3 adds
 * doubles and booleans, so every scalar shape is handled here rather than at
 * each call site.
 */
func redisValueToString(value interface{}) string {
	text, _ := redisValueToStringOK(value)
	return text
}

/**
 * redisValueToStringOK renders a RESP scalar as text, reporting whether the
 * value was a scalar at all. Arrays, maps and unknown types return false so a
 * caller can skip them instead of inserting an empty string.
 */
func redisValueToStringOK(value interface{}) (string, bool) {
	switch value := value.(type) {
	case []byte:
		return string(value), true
	case string:
		return value, true
	case int64:
		return strconv.FormatInt(value, 10), true
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64), true
	case bool:
		return strconv.FormatBool(value), true
	default:
		return "", false
	}
}

/**
 * redisValueToInt64 converts a RESP value to an integer, reporting whether the
 * conversion was possible.
 */
func redisValueToInt64(value interface{}) (int64, bool) {
	switch value := value.(type) {
	case int64:
		return value, true
	case float64:
		return int64(value), true
	case []byte:
		parsed, err := strconv.ParseInt(string(value), 10, 64)
		return parsed, err == nil
	case string:
		parsed, err := strconv.ParseInt(value, 10, 64)
		return parsed, err == nil
	default:
		return 0, false
	}
}
