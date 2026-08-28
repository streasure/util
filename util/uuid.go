package util

import (
	"github.com/google/uuid"
	"strings"
)

func NewUUID() string {
	u4 := uuid.New()
	return strings.ReplaceAll(u4.String(), "-", "")
}

func NewUUIDBytes() []byte {
	u4 := uuid.New()
	return u4[:]
}
