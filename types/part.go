package types

import "github.com/gotd/td/tg"

type Part struct {
	Location *tg.InputDocumentFileLocation
	Start    int64
	End      int64
}
