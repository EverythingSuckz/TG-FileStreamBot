// This file is a part of EverythingSuckz/TG-FileStreamBot
// And is licenced under the Affero General Public License.
// Any distributions of this code MUST be accompanied by a copy of the AGPL
// with proper attribution to the original author(s).

package qrlogin

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"strings"

	"github.com/gotd/td/session"
)

func EncodeToPyrogramSession(data *session.Data, appID int32) (string, error) {
	buf := new(bytes.Buffer)
	if err := buf.WriteByte(byte(data.DC)); err != nil {
		return "", err
	}
	if err := binary.Write(buf, binary.BigEndian, appID); err != nil {
		return "", err
	}
	var testMode byte
	if data.Config.TestMode {
		testMode = 1
	}
	if err := buf.WriteByte(testMode); err != nil {
		return "", err
	}
	if len(data.AuthKey) != 256 {
		return "", errors.New("auth key must be 256 bytes long")
	}
	if _, err := buf.Write(data.AuthKey); err != nil {
		return "", err
	}
	if len(data.AuthKeyID) != 8 {
		return "", errors.New("auth key ID must be 8 bytes long")
	}
	if _, err := buf.Write(data.AuthKeyID); err != nil {
		return "", err
	}
	if err := buf.WriteByte(0); err != nil {
		return "", err
	}
	// Convert the bytes buffer to a base64 string
	encodedString := base64.URLEncoding.EncodeToString(buf.Bytes())
	trimmedEncoded := strings.TrimRight(encodedString, "=")
	return trimmedEncoded, nil
}
