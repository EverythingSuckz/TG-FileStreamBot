package utils

import (
	"context"
	"sync"
	"time"
)

// ChatPacer keeps at least a fixed delay between operations that belong to the
// same chat.
//
// A media group is delivered as one update per message, so a single album of 8
// photos makes the bot forward 8 messages and send 8 replies in a couple of
// seconds. Telegram answers that burst with 420 FLOOD_WAIT. The retry
// middleware absorbs those answers, this pacer keeps most of them from
// happening in the first place.
type ChatPacer struct {
	interval time.Duration

	mu   sync.Mutex
	next map[int64]time.Time
}

// NewChatPacer returns a pacer that allows one operation per interval for a
// given chat. A non-positive interval disables pacing.
func NewChatPacer(interval time.Duration) *ChatPacer {
	return &ChatPacer{
		interval: interval,
		next:     make(map[int64]time.Time),
	}
}

// Wait blocks until the next operation for chatID is allowed to run.
func (p *ChatPacer) Wait(ctx context.Context, chatID int64) error {
	if p == nil || p.interval <= 0 {
		return nil
	}
	now := time.Now()

	p.mu.Lock()
	if len(p.next) > 4096 {
		for id, at := range p.next {
			if at.Before(now) {
				delete(p.next, id)
			}
		}
	}
	at := p.next[chatID]
	if at.Before(now) {
		at = now
	}
	p.next[chatID] = at.Add(p.interval)
	p.mu.Unlock()

	delay := time.Until(at)
	if delay <= 0 {
		return nil
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// SendPacer is the process wide pacer used before answering a user.
var SendPacer = NewChatPacer(0)

// InitSendPacer configures SendPacer from the loaded configuration.
func InitSendPacer(interval time.Duration) {
	SendPacer = NewChatPacer(interval)
}
