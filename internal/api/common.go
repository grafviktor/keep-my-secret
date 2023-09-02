package api

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ContextKey string

const ContextUserLogin ContextKey = "login"
