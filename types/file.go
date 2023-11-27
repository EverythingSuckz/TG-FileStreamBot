package types

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"fmt"

	"github.com/gotd/td/tg"
)

type File struct {
	Location *tg.InputDocumentFileLocation
	FileSize int64
	FileName string
	MimeType string
	ID       int64
}

type HashableFileStruct struct {
	FileName string
	FileSize int64
	MimeType string
	FileID   int64
}

func (f *HashableFileStruct) Pack() string {
	var b bytes.Buffer
	gob.NewEncoder(&b).Encode(f)
	return FromBytes(b.Bytes())
}

func FromBytes(data []byte) string {
	result := md5.Sum(data)
	return fmt.Sprintf("%x", result)
}
