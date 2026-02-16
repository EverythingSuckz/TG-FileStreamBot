package stream

import (
	"EverythingSuckz/fsb/config"
	"EverythingSuckz/fsb/internal/utils"
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/celestix/gotgproto"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

// calculateBlockSize func determines optimal block size based on the range requested.
// Smaller ranges use smaller blocks to reduce wasted bandwidth during seeks.
func calculateBlockSize(start, end int64) int64 {
	size := end - start + 1

	switch {
	case size < 512*1024: // < 512KB
		return 64 * 1024 // 64KB blocks
	case size < 4*1024*1024: // < 4MB
		return 256 * 1024 // 256KB blocks
	case size < 32*1024*1024: // < 32MB
		return 512 * 1024 // 512KB blocks
	default:
		return 1024 * 1024 // 1MB blocks by default
	}
}

// StreamPipe reads data from Telegram with concurrent prefetching. implements `io.ReadCloser`.
// credits: [teldrive](https://github.com/tgdrive/teldrive/blob/1071a6e8b5a4076cacc5a0989ef245fae517d837/internal/reader/tg_reader.go)
type StreamPipe struct {
	ctx    context.Context
	cancel context.CancelFunc
	log    *zap.Logger

	// tg connection
	client   *gotgproto.Client
	location tg.InputFileLocationClass

	// range stuff
	start      int64
	end        int64
	blockSize  int64
	totalBytes int64

	// prefetch pipeline
	blockQueue chan []byte

	// current read state
	currentBlock []byte
	blockOffset  int64
	bytesRead    int64

	// lifecycle
	closeOnce sync.Once
}

// class function which creates a StreamPipe with default configuration.
func NewStreamPipe(
	ctx context.Context,
	client *gotgproto.Client,
	location tg.InputFileLocationClass,
	start, end int64,
	log *zap.Logger,
) (io.ReadCloser, error) {
	ctx, cancel := context.WithCancel(ctx)

	totalBytes := end - start + 1
	blockSize := calculateBlockSize(start, end)

	p := &StreamPipe{
		ctx:        ctx,
		cancel:     cancel,
		log:        log.Named("StreamPipe"),
		client:     client,
		location:   location,
		start:      start,
		end:        end,
		blockSize:  blockSize,
		totalBytes: totalBytes,
		blockQueue: make(chan []byte, config.ValueOf.StreamBufferCount),
	}

	// start prefetching in background
	go p.prefetch()

	return p, nil
}

// Read implements io.Reader
func (p *StreamPipe) Read(buf []byte) (n int, err error) {
	if p.bytesRead >= p.totalBytes {
		p.log.Sugar().Debug("EOF (bytesread == contentLength)")
		return 0, io.EOF
	}

	// need a new block?
	if p.blockOffset >= int64(len(p.currentBlock)) {
		select {
		case block, ok := <-p.blockQueue:
			if !ok {
				if p.bytesRead >= p.totalBytes {
					return 0, io.EOF
				}
				return 0, ErrPipeDrained
			}
			p.currentBlock = block
			p.blockOffset = 0
		case <-p.ctx.Done():
			return 0, p.ctx.Err()
		}
	}

	// copy available data
	n = copy(buf, p.currentBlock[p.blockOffset:])
	p.blockOffset += int64(n)
	p.bytesRead += int64(n)

	return n, nil
}

// Close implements io.Closer.
// it cancels prefetching and releases resources.
func (p *StreamPipe) Close() error {
	p.closeOnce.Do(func() {
		p.cancel()
	})
	return nil
}

// prefetch runs in a goroutine, fetching blocks concurrently and sending to blockQueue.
func (p *StreamPipe) prefetch() {
	defer close(p.blockQueue)

	// calc block boundaries
	alignedStart := p.start - (p.start % p.blockSize)
	leftTrim := p.start - alignedStart
	rightTrim := (p.end % p.blockSize) + 1
	totalBlocks := int((p.end - alignedStart + p.blockSize) / p.blockSize)

	currentBlock := 0
	offset := alignedStart

	for currentBlock < totalBlocks {
		// check for cancellation
		select {
		case <-p.ctx.Done():
			return
		default:
		}

		// fetch a batch of blocks concurrently
		batchSize := min(config.ValueOf.StreamConcurrency, totalBlocks-currentBlock)
		blocks := make([][]byte, batchSize)

		var wg sync.WaitGroup
		var fetchErr error
		var errMu sync.Mutex

		for i := range batchSize {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()

				blockNum := currentBlock + idx
				blockOffset := offset + int64(idx)*p.blockSize

				data, err := p.downloadBlockWithRetry(blockOffset)
				dataLen := int64(len(data))

				if err != nil {
					errMu.Lock()
					if fetchErr == nil {
						fetchErr = err
					}
					errMu.Unlock()
					return
				}

				// trim first/last block to exact range
				if totalBlocks == 1 {
					if dataLen < rightTrim {
						rightTrim = dataLen
					}
					if leftTrim > dataLen {
						leftTrim = dataLen
					}
					data = data[leftTrim:rightTrim]
				} else if blockNum == 0 {
					if leftTrim > dataLen {
						leftTrim = dataLen
					}
					data = data[leftTrim:]
				} else if blockNum == totalBlocks-1 {
					if dataLen > rightTrim {
						data = data[:rightTrim]
					}
				}

				blocks[idx] = data
			}(i)
		}

		wg.Wait()

		// handle errors
		// ignore context cancellation cuz it's expected on disconnect
		if fetchErr != nil {
			if p.ctx.Err() == nil {
				p.log.Error("block download failed", zap.Error(fetchErr))
			}
			return
		}

		// send blocks to queue in order
		for _, block := range blocks {
			if block == nil {
				// a fetch failure that wasn't captured, should not happen but just in case.
				p.log.Error("unexpected nil block in batch, aborting prefetch")
				return
			}
			select {
			case p.blockQueue <- block:
			case <-p.ctx.Done():
				return
			}
		}

		currentBlock += batchSize
		offset += p.blockSize * int64(batchSize)
	}
}

// downloadBlockWithRetry fetches a block with exponential backoff retry.
func (p *StreamPipe) downloadBlockWithRetry(offset int64) ([]byte, error) {
	var lastErr error

	// TODO: make configurable later
	backoff := 100 * time.Millisecond   // initial backoff = 100ms
	const maxBackoff = 15 * time.Second // max backoff = 15s

	for attempt := 0; attempt < config.ValueOf.StreamMaxRetries; attempt++ {
		// check context before each attempt
		if p.ctx.Err() != nil {
			return nil, p.ctx.Err()
		}

		ctx, cancel := context.WithTimeout(p.ctx, time.Duration(config.ValueOf.StreamTimeoutSec)*time.Second)
		data, err := utils.TimeFuncWithResult(p.log, "downloadBlock", func() ([]byte, error) {
			return p.downloadBlock(ctx, offset)
		})
		cancel()

		if err == nil {
			return data, nil
		}

		lastErr = err

		// don't retry on context cancellation
		if p.ctx.Err() != nil {
			return nil, p.ctx.Err()
		}

		// exponential backoff
		select {
		case <-time.After(backoff):
			backoff *= 2
			// making sure backoff doesn't grow indefinitely in case of persistent failures
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		case <-p.ctx.Done():
			return nil, p.ctx.Err()
		}
	}

	return nil, fmt.Errorf("%w: %v", ErrMaxRetriesExceeded, lastErr)
}

// downloadBlock fetches a single block from Telegram.
func (p *StreamPipe) downloadBlock(ctx context.Context, offset int64) ([]byte, error) {
	p.log.Sugar().Debugf("Downloading block at offset %d (block size: %d)", offset, p.blockSize)
	res, err := p.client.API().UploadGetFile(ctx, &tg.UploadGetFileRequest{
		Offset:   offset,
		Limit:    int(p.blockSize),
		Location: p.location,
	})

	if err != nil {
		return nil, err
	}

	switch result := res.(type) {
	case *tg.UploadFile:
		return result.Bytes, nil
	default:
		return nil, fmt.Errorf("unexpected response type: %T", res)
	}
}
