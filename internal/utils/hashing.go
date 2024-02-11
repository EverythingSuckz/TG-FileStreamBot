package utils

import (
	"EverythingSuckz/fsb/config"
	"EverythingSuckz/fsb/internal/types"
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
