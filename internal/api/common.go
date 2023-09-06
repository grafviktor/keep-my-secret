package api

// Response - is a struct which unifies JSON response format for the client
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// ContextKey - is a type for http request context keys
type ContextKey string

// ContextUserLogin - used to store user login in context
const ContextUserLogin ContextKey = "login"
