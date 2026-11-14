package bot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"EverythingSuckz/fsb/config"

	"github.com/gotd/contrib/middleware/floodwait"
	"github.com/gotd/contrib/middleware/ratelimit"
	"github.com/gotd/td/bin"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgerr"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// floodWaitLogger reports FLOOD_WAIT answers from Telegram. It sits between the
// waiter and the rate limiter, so every single attempt that got throttled is
// logged, not just the final failure.
type floodWaitLogger struct {
	log *zap.Logger
}

func (m floodWaitLogger) Handle(next tg.Invoker) telegram.InvokeFunc {
	return func(ctx context.Context, input bin.Encoder, output bin.Decoder) error {
		err := next.Invoke(ctx, input, output)
		if d, ok := tgerr.AsFloodWait(err); ok {
			m.log.Warn("Telegram asked to wait, retrying",
				zap.String("method", methodName(input)),
				zap.Duration("wait", d),
			)
		}
		return err
	}
}

// methodName strips the generated type decoration, e.g.
// *tg.MessagesForwardMessagesRequest becomes MessagesForwardMessages.
func methodName(input bin.Encoder) string {
	name := fmt.Sprintf("%T", input)
	name = strings.TrimPrefix(name, "*tg.")
	name = strings.TrimSuffix(name, "Request")
	return name
}

// GetFloodMiddleware returns the middleware chain used by every Telegram
// client (main bot, workers and userbot).
//
// Order matters: middlewares are applied in slice order, so the waiter is the
// outermost one and its retries go through the logger and the rate limiter
// again. That way a retried request is both logged and rate limited.
func GetFloodMiddleware(log *zap.Logger) []telegram.Middleware {
	return newFloodMiddleware(
		log.Named("FloodWait"),
		uint(config.ValueOf.FloodWaitMaxRetries),
		time.Duration(config.ValueOf.FloodWaitMaxWaitSec)*time.Second,
		time.Duration(config.ValueOf.RateLimitIntervalMs)*time.Millisecond,
		config.ValueOf.RateLimitBurst,
	)
}

// newFloodMiddleware builds the chain from explicit values so it can be
// exercised in tests without touching the global config.
func newFloodMiddleware(
	log *zap.Logger,
	maxRetries uint,
	maxWait time.Duration,
	interval time.Duration,
	burst int,
) []telegram.Middleware {
	waiter := floodwait.NewSimpleWaiter().
		WithMaxRetries(maxRetries).
		WithMaxWait(maxWait)
	limiter := ratelimit.New(rate.Every(interval), burst)
	return []telegram.Middleware{
		waiter,
		floodWaitLogger{log: log},
		limiter,
	}
}
