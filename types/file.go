package types

import "github.com/gotd/td/tg"

type File struct {
	Location *tg.InputDocumentFileLocation
	FileSize int64
	FileName string
	MimeType string
	ID       int64
}
