package trace

import "context"

type contextKey string

const (
	// ID is the key to store the trace ID in the context
	ID contextKey = "traceID"
)

// GetID returns the trace ID from the context
func GetID(ctx context.Context) (string, string) {
	if ctx == nil || ctx.Value(ID) == nil {
		return "", ""
	}
	return string(ID), ctx.Value(ID).(string)
}

// String returns the string representation of the context key
func (key contextKey) String() string {
	return string(key)
}
