package keybuilder

import (
	"fmt"
)

const (
	// A common prefix for all keys related to this service.
	redisPrefix = "shortener"
	// The entity type we are caching.
	urlKey = "url"
)

// URLCacheKey builds a standardized Redis key for a URL cache entry.
func URLCacheKey(shortCode string) string {
	return fmt.Sprintf("%s:%s:%s", redisPrefix, urlKey, shortCode)
}
