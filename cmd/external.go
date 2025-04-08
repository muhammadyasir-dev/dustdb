package main

import (
	"errors"
	"strconv"
	"time"
)

// Common errors for numeric operations
var (
	ErrNotNumeric    = errors.New("value is not numeric")
	ErrOverflow      = errors.New("operation would result in overflow")
	ErrUnderflow     = errors.New("operation would result in underflow")
	ErrKeyNotFound   = errors.New("key not found")
	ErrInvalidAmount = errors.New("invalid increment/decrement amount")
)

// NumericOperations provides additional numeric operations for the cache
type NumericOperations struct {
	cache *Cache
}

// NewNumericOperations creates a new NumericOperations instance
func NewNumericOperations(cache *Cache) *NumericOperations {
	return &NumericOperations{
		cache: cache,
	}
}

// INCR increments the value of a key by an amount
// If the key does not exist, it's initialized with 0 before incrementing
// Returns the new value after incrementing or an error
func (n *NumericOperations) INCR(key string, amount int64) (int64, error) {
	if amount == 0 {
		return 0, ErrInvalidAmount
	}

	shard := n.cache.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	// Get current value or initialize with 0
	var currentValue int64 = 0
	entry, exists := shard.data[key]

	if exists {
		// Check if the value is a valid number
		val, err := strconv.ParseInt(entry.Value, 10, 64)
		if err != nil {
			return 0, ErrNotNumeric
		}
		currentValue = val
	}

	// Check for potential overflow
	if (amount > 0 && currentValue > (1<<63-1)-amount) ||
		(amount < 0 && currentValue < (-1<<63)+amount) {
		return 0, ErrOverflow
	}

	// Calculate new value
	newValue := currentValue + amount

	// Update the cache
	shard.data[key] = CacheEntry{
		Value:    strconv.FormatInt(newValue, 10),
		ExpireAt: entry.ExpireAt, // Preserve TTL if it exists
	}

	return newValue, nil
}

// DECR decrements the value of a key by an amount
// If the key does not exist, it's initialized with 0 before decrementing
// Returns the new value after decrementing or an error
func (n *NumericOperations) DECR(key string, amount int64) (int64, error) {
	// Reuse INCR logic with negative amount
	if amount <= 0 {
		return 0, ErrInvalidAmount
	}
	return n.INCR(key, -amount)
}

// IncrBy increments the value of a key by the specified amount
// Returns the new value or an error
func (n *NumericOperations) IncrBy(key string, amount int64) (int64, error) {
	return n.INCR(key, amount)
}

// DecrBy decrements the value of a key by the specified amount
// Returns the new value or an error
func (n *NumericOperations) DecrBy(key string, amount int64) (int64, error) {
	return n.DECR(key, amount)
}

// Incr increments the value of a key by 1
// Returns the new value or an error
func (n *NumericOperations) Incr(key string) (int64, error) {
	return n.INCR(key, 1)
}

// Decr decrements the value of a key by 1
// Returns the new value or an error
func (n *NumericOperations) Decr(key string) (int64, error) {
	return n.DECR(key, 1)
}

// UpdateExpiration updates the expiration time of a key
// Returns true if the key exists and its expiration was updated
func (n *NumericOperations) UpdateExpiration(key string, ttl time.Duration) bool {
	shard := n.cache.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	entry, exists := shard.data[key]
	if !exists {
		return false
	}

	var expireAt int64
	if ttl > 0 {
		expireAt = time.Now().Add(ttl).UnixNano()
	}

	entry.ExpireAt = expireAt
	shard.data[key] = entry

	return true
}
