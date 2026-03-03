package utils

import (
	"EverythingSuckz/fsb/config"
	"EverythingSuckz/fsb/internal/types"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

func PackFile(fileName string, fileSize int64, mimeType string, fileID int64) string {
	return (&types.HashableFileStruct{FileName: fileName, FileSize: fileSize, MimeType: mimeType, FileID: fileID}).Pack()
}

func GetShortHash(fullHash string) string {
	return fullHash[:config.ValueOf.HashLength]
}

func CheckHash(inputHash string, expectedHash string) bool {
	return inputHash == GetShortHash(expectedHash)
}

// SignURL generates a truncated HMAC-SHA256 signature for messageID:expiry.
func SignURL(messageID int, expiry int64) string {
	payload := fmt.Sprintf("%d:%d", messageID, expiry)
	mac := hmac.New(sha256.New, []byte(config.ValueOf.BotToken))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))[:config.ValueOf.HashLength]
}

func VerifyURL(sig string, messageID int, expiryStr string) (string, bool) {
	expiry, err := strconv.ParseInt(expiryStr, 10, 64)
	if err != nil {
		return "invalid expiry", false
	}
	if expiry != 0 && time.Now().Unix() > expiry {
		return "link has expired", false
	}
	expected := SignURL(messageID, expiry)
	if !hmac.Equal([]byte(sig), []byte(expected)) {
		return "invalid signature", false
	}
	return "", true
}
