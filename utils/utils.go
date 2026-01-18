package utils

import (
	"time"

	"github.com/google/uuid"
)

// GenerateUUID generates a new UUID
func GenerateUUID() string {
	return uuid.New().String()
}

// TimeNowUTC returns current time in UTC
func TimeNowUTC() time.Time {
	return time.Now().UTC()
}

// Pointer returns a pointer to the given value
func Pointer[T any](v T) *T {
	return &v
}

// StringPtrToString returns string value from pointer, or empty string if nil
func StringPtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// IntPtrToInt returns int value from pointer, or 0 if nil
func IntPtrToInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}
