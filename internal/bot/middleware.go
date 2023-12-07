package bot

import (
	"time"

	"github.com/gotd/contrib/middleware/floodwait"
	"github.com/gotd/contrib/middleware/ratelimit"
	"github.com/gotd/td/telegram"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

func GetFloodMiddleware(log *zap.Logger) []telegram.Middleware {
	waiter := floodwait.NewSimpleWaiter().WithMaxRetries(10)
	ratelimiter := ratelimit.New(rate.Every(time.Millisecond*100), 5)
	return []telegram.Middleware{
		waiter,
		ratelimiter,
	}
}
