package bot

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gotd/td/bin"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgerr"
	"go.uber.org/zap"
)

func chain(middlewares []telegram.Middleware, next tg.Invoker) tg.Invoker {
	for i := len(middlewares) - 1; i >= 0; i-- {
		next = middlewares[i].Handle(next)
	}
	return next
}

func TestFloodMiddlewareRetriesFloodWait(t *testing.T) {
	var attempts atomic.Int32
	invoker := telegram.InvokeFunc(func(ctx context.Context, input bin.Encoder, output bin.Decoder) error {
		if attempts.Add(1) < 2 {
			return tgerr.New(420, "FLOOD_WAIT_0")
		}
		return nil
	})

	mw := newFloodMiddleware(zap.NewNop(), 5, time.Minute, time.Millisecond, 1)
	if err := chain(mw, invoker).Invoke(context.Background(), &tg.HelpGetConfigRequest{}, nil); err != nil {
		t.Fatalf("invoke returned error: %v", err)
	}
	if got := attempts.Load(); got != 2 {
		t.Fatalf("attempts = %d, want 2", got)
	}
}

func TestFloodMiddlewareHonoursMaxWait(t *testing.T) {
	var attempts atomic.Int32
	invoker := telegram.InvokeFunc(func(ctx context.Context, input bin.Encoder, output bin.Decoder) error {
		attempts.Add(1)
		return tgerr.New(420, "FLOOD_WAIT_600")
	})

	mw := newFloodMiddleware(zap.NewNop(), 5, 10*time.Second, time.Millisecond, 1)
	err := chain(mw, invoker).Invoke(context.Background(), &tg.HelpGetConfigRequest{}, nil)
	if err == nil {
		t.Fatal("invoke returned nil, want error")
	}
	if _, ok := tgerr.AsFloodWait(err); !ok {
		t.Fatalf("error %v does not wrap the flood wait", err)
	}
	if got := attempts.Load(); got != 1 {
		t.Fatalf("attempts = %d, want 1", got)
	}
}

func TestMethodName(t *testing.T) {
	if got := methodName(&tg.MessagesForwardMessagesRequest{}); got != "MessagesForwardMessages" {
		t.Fatalf("methodName = %q", got)
	}
}
