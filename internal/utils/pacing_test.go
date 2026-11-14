package utils

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestChatPacerConcurrentSameChat(t *testing.T) {
	p := NewChatPacer(10 * time.Millisecond)
	start := time.Now()
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			if err := p.Wait(context.Background(), 1); err != nil {
				t.Error(err)
			}
		})
	}
	wg.Wait()
	if d := time.Since(start); d < 70*time.Millisecond {
		t.Fatalf("8 concurrent waits took %v, want at least 70ms", d)
	}
}

func TestChatPacerSpacesSameChat(t *testing.T) {
	p := NewChatPacer(20 * time.Millisecond)
	start := time.Now()
	for i := range 4 {
		if err := p.Wait(context.Background(), 1); err != nil {
			t.Fatalf("wait %d: %v", i, err)
		}
	}
	if d := time.Since(start); d < 60*time.Millisecond {
		t.Fatalf("4 waits took %v, want at least 60ms", d)
	}
}

func TestChatPacerSeparatesChats(t *testing.T) {
	p := NewChatPacer(time.Hour)
	if err := p.Wait(context.Background(), 1); err != nil {
		t.Fatalf("first wait: %v", err)
	}
	start := time.Now()
	if err := p.Wait(context.Background(), 2); err != nil {
		t.Fatalf("second wait: %v", err)
	}
	if d := time.Since(start); d > time.Second {
		t.Fatalf("other chat was delayed by %v", d)
	}
}

func TestChatPacerStopsOnCanceledContext(t *testing.T) {
	p := NewChatPacer(time.Hour)
	if err := p.Wait(context.Background(), 1); err != nil {
		t.Fatalf("first wait: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := p.Wait(ctx, 1); err == nil {
		t.Fatal("wait returned nil on canceled context")
	}
}

func TestChatPacerDisabled(t *testing.T) {
	p := NewChatPacer(0)
	start := time.Now()
	for i := range 100 {
		if err := p.Wait(context.Background(), 1); err != nil {
			t.Fatalf("wait %d: %v", i, err)
		}
	}
	if d := time.Since(start); d > time.Second {
		t.Fatalf("disabled pacer blocked for %v", d)
	}
}
