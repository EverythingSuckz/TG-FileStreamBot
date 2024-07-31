package bot

import (
	"EverythingSuckz/fsb/config"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/sessionMaker"
	"github.com/glebarez/sqlite"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

type Worker struct {
	ID     int
	Client *gotgproto.Client
	Self   *tg.User
	log    *zap.Logger
}

func (w *Worker) String() string {
	return fmt.Sprintf("{Worker (%d|@%s)}", w.ID, w.Self.Username)
}

type BotWorkers struct {
	Bots     []*Worker
	starting int
	index    int
	mut      sync.Mutex
	log      *zap.Logger
}

var Workers *BotWorkers = &BotWorkers{
	log:  nil,
	Bots: make([]*Worker, 0),
}

func (w *BotWorkers) Init(log *zap.Logger) {
	w.log = log.Named("Workers")
}

func (w *BotWorkers) AddDefaultClient(client *gotgproto.Client, self *tg.User) {
	if w.Bots == nil {
		w.Bots = make([]*Worker, 0)
	}
	w.incStarting()
	w.Bots = append(w.Bots, &Worker{
		Client: client,
		ID:     w.starting,
		Self:   self,
		log:    w.log,
	})
	w.log.Sugar().Info("Default bot loaded")
}

func (w *BotWorkers) incStarting() {
	w.mut.Lock()
	defer w.mut.Unlock()
	w.starting++
}

func (w *BotWorkers) Add(token string) (err error) {
	w.incStarting()
	var botID int = w.starting
	client, err := startWorker(w.log, token, botID)
	if err != nil {
		return err
	}
	w.log.Sugar().Infof("Bot @%s loaded with ID %d", client.Self.Username, botID)
	w.Bots = append(w.Bots, &Worker{
		Client: client,
		ID:     botID,
		Self:   client.Self,
		log:    w.log,
	})
	return nil
}

func GetNextWorker() *Worker {
	Workers.mut.Lock()
	defer Workers.mut.Unlock()
	index := (Workers.index + 1) % len(Workers.Bots)
	Workers.index = index
	worker := Workers.Bots[index]
	Workers.log.Sugar().Debugf("Using worker %d", worker.ID)
	return worker
}

func StartWorkers(log *zap.Logger) (*BotWorkers, error) {
	Workers.Init(log)

	if len(config.ValueOf.MultiTokens) == 0 {
		Workers.log.Sugar().Info("No worker bot tokens provided, skipping worker initialization")
		return Workers, nil
	}
	Workers.log.Sugar().Info("Starting")
	if config.ValueOf.UseSessionFile {
		Workers.log.Sugar().Info("Using session file for workers")
		newpath := filepath.Join(".", "sessions")
		if err := os.MkdirAll(newpath, os.ModePerm); err != nil {
			Workers.log.Error("Failed to create sessions directory", zap.Error(err))
			return nil, err
		}
	}

	var wg sync.WaitGroup
	var successfulStarts int32
	totalBots := len(config.ValueOf.MultiTokens)

	for i := 0; i < totalBots; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			done := make(chan error, 1)
			go func() {
				err := Workers.Add(config.ValueOf.MultiTokens[i])
				done <- err
			}()

			select {
			case err := <-done:
				if err != nil {
					Workers.log.Error("Failed to start worker", zap.Int("index", i), zap.Error(err))
				} else {
					atomic.AddInt32(&successfulStarts, 1)
				}
			case <-ctx.Done():
				Workers.log.Error("Timed out starting worker", zap.Int("index", i))
			}
		}(i)
	}

	wg.Wait() // Wait for all goroutines to finish
	Workers.log.Sugar().Infof("Successfully started %d/%d bots", successfulStarts, totalBots)
	return Workers, nil
}

func startWorker(l *zap.Logger, botToken string, index int) (*gotgproto.Client, error) {
	log := l.Named("Worker").Sugar()
	log.Infof("Starting worker with index - %d", index)
	var sessionType sessionMaker.SessionConstructor
	if config.ValueOf.UseSessionFile {
		sessionType = sessionMaker.SqlSession(sqlite.Open(fmt.Sprintf("sessions/worker-%d.session", index)))
	} else {
		sessionType = sessionMaker.SimpleSession()
	}
	client, err := gotgproto.NewClient(
		int(config.ValueOf.ApiID),
		config.ValueOf.ApiHash,
		gotgproto.ClientTypeBot(botToken),
		&gotgproto.ClientOpts{
			Session:          sessionType,
			DisableCopyright: true,
			Middlewares:      GetFloodMiddleware(log.Desugar()),
		},
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}
