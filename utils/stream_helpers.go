package utils

import (
	"EverythingSuckz/fsb/types"
	"bytes"
	"context"
	"fmt"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

type Parts []types.Part

func getChunk(ctx context.Context, tgClient *telegram.Client, location tg.InputFileLocationClass, offset int64, limit int64) ([]byte, error) {
	req := &tg.UploadGetFileRequest{
		Offset:   offset,
		Limit:    int(limit),
		Location: location,
	}
	r, err := tgClient.API().UploadGetFile(ctx, req)
	if err != nil {
		return nil, err
	}
	switch result := r.(type) {
	case *tg.UploadFile:
		return result.Bytes, nil
	default:
		return nil, fmt.Errorf("unexpected type %T", r)
	}
}

func IterContent(ctx context.Context, tgClient *telegram.Client, location tg.InputFileLocationClass) (*bytes.Buffer, error) {
	offset := int64(0)
	limit := int64(1024 * 1024)
	buff := &bytes.Buffer{}
	for {
		r, err := getChunk(ctx, tgClient, location, offset, limit)
		if err != nil {
			return buff, err
		}
		if len(r) == 0 {
			break
		}
		buff.Write(r)
		offset += int64(limit)
	}
	return buff, nil
}

func GetParts(ctx context.Context, client *telegram.Client, file *types.File) ([]types.Part, error) {
	parts := []types.Part{}
	parts = append(parts, types.Part{Location: file.Location, Start: 0, End: file.FileSize - 1})
	return parts, nil
}

func RangedParts(parts []types.Part, startByte, endByte int64) []types.Part {
	chunkSize := parts[0].End + 1
	numParts := int64(len(parts))
	validParts := []types.Part{}
	firstChunk := max(startByte/chunkSize, 0)
	lastChunk := min(endByte/chunkSize, numParts)
	startInFirstChunk := startByte % chunkSize
	endInLastChunk := endByte % chunkSize
	if firstChunk == lastChunk {
		validParts = append(validParts, types.Part{
			Location: parts[firstChunk].Location,
			Start:    startInFirstChunk,
			End:      endInLastChunk,
		})
	} else {
		validParts = append(validParts, types.Part{
			Location: parts[firstChunk].Location,
			Start:    startInFirstChunk,
			End:      parts[firstChunk].End,
		})
		// Add valid parts from any chunks in between.
		for i := firstChunk + 1; i < lastChunk; i++ {
			validParts = append(validParts, types.Part{
				Location: parts[i].Location,
				Start:    0,
				End:      parts[i].End,
			})
		}
		// Add valid parts from the last chunk.
		validParts = append(validParts, types.Part{
			Location: parts[lastChunk].Location,
			Start:    0,
			End:      endInLastChunk,
		})
	}
	return validParts
}
