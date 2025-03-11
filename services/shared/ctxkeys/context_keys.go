package ctxkeys

type contextKey string

const (
	UsernameKey  contextKey = "username"
	RequestIDKey contextKey = "request_id"
)
