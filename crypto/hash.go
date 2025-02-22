package crypto

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"errors"
	"github.com/google/uuid"
	"github.com/zerok-ai/zk-utils-go/common"
	"github.com/zerok-ai/zk-utils-go/logs"
	"io"
)

var LogTag = "zk_hash"

func CalculateHashNewSHA2(s string) uuid.UUID {
	hash := sha1.Sum([]byte(s))
	return uuid.NewSHA1(uuid.Nil, hash[:])
}

func CompressStringGzip(input string) ([]byte, error) {
	if common.IsEmpty(input) {
		logger.Error(LogTag, "empty string")
		return nil, errors.New("empty input")
	}
	var b bytes.Buffer

	gzipWriter := gzip.NewWriter(&b)

	_, err := gzipWriter.Write([]byte(input))
	if err != nil {
		logger.Error(LogTag, err)
		return nil, err
	}

	err = gzipWriter.Close()
	if err != nil {
		logger.Error(LogTag, err)
		return nil, err
	}

	return b.Bytes(), nil
}

func DecompressStringGzip(input []byte) (string, error) {
	if input == nil || len(input) == 0 {
		return "", errors.New("empty input")
	}

	reader := bytes.NewReader(input)

	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		logger.Error(LogTag, err)
		return "", err
	}
	defer gzipReader.Close()

	decompressed, err := io.ReadAll(gzipReader)
	if err != nil {
		logger.Error(LogTag, err)
		return "", err
	}

	logger.Debug(LogTag, "String decompressed successfully")

	return string(decompressed), nil
}
