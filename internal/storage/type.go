// Package storage defines the storage types
package storage

// Type - storage type. Used to determine which storage is used by the application
type Type string // Storage type

// Storage types
const (
	TypeMemory Type = "memoryStorage" // Memory storage, when other storage types aren't available
	TypeSQL    Type = "sqlStorage"    // SQL storage
)
