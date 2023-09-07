// Package storage defines the storage types
package storage

// Type - storage type. Used to determine which storage is used by the application
type Type string // Storage type

// TypeSQL Storage - only sql storage is supported at the moment
const (
	TypeSQL Type = "sqlStorage" // SQL storage
)
