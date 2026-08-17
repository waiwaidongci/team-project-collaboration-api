package taskrules

import "context"

func WithTaskContext(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(parent)
}
