package bot

import (
	"EverythingSuckz/fsb/config"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/sessionMaker"
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
	Workers.log.Sugar().Infof("Using worker %d", worker.ID)
	return worker
}

func StartWorkers(log *zap.Logger) {
	log.Sugar().Info("Starting workers")
	Workers.Init(log)
	if config.ValueOf.UseSessionFile {
		log.Sugar().Info("Using session file for workers")
		newpath := filepath.Join(".", "sessions")
		err := os.MkdirAll(newpath, os.ModePerm)
		if err != nil {
			log.Error("Failed to create sessions directory", zap.Error(err))
			return
		}
	}
	c := make(chan struct{})
	for i := 0; i < len(config.ValueOf.MultiTokens); i++ {
		go func(i int) {
			err := Workers.Add(config.ValueOf.MultiTokens[i])
			if err != nil {
				log.Error("Failed to start worker", zap.Error(err))
				return
			}
			c <- struct{}{}
		}(i)
	}
	// wait for all workers to start
	log.Sugar().Info("Waiting for all workers to start")
	for i := 0; i < len(config.ValueOf.MultiTokens); i++ {
		<-c
	}
}

func startWorker(l *zap.Logger, botToken string, index int) (*gotgproto.Client, error) {
	log := l.Named("Worker").Sugar()
	log.Infof("Starting worker with index - %d", index)
	var sessionType *sessionMaker.SqliteSessionConstructor
	if config.ValueOf.UseSessionFile {
		sessionType = sessionMaker.SqliteSession(fmt.Sprintf("sessions/worker-%d", index))
	} else {
		sessionType = sessionMaker.SqliteSession(":memory:")
	}
	client, err := gotgproto.NewClient(
		int(config.ValueOf.ApiID),
		config.ValueOf.ApiHash,
		gotgproto.ClientType{
			BotToken: botToken,
		},
		&gotgproto.ClientOpts{
			Session:          sessionType,
			DisableCopyright: true,
		},
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}
