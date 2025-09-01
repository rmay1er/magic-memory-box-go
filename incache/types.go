package incache

import (
	"time"
)

type mapFields struct {
	Value      any
	ExpireTime time.Time
}

type inCache struct {
	memory map[string]mapFields
}
