package timeout

import (
	"time"

	context "golang.org/x/net/context"
)

// AuthContext provides a context with a timeout appropriate for auth
func AuthContext() context.Context {
	clientDeadline := time.Now().Add(time.Duration(time.Millisecond * 3000))
	ctx, _ := context.WithDeadline(context.Background(), clientDeadline)

	return ctx
}
