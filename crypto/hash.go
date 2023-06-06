package crypto

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"errors"
	"github.com/google/uuid"
	zkcommon "github.com/zerok-ai/zk-utils-go/common"
	"io"
	"log"
)

func CalculateHash(s string) uuid.UUID {
	hash := sha1.Sum([]byte(s))
	return uuid.NewSHA1(uuid.Nil, hash[:])
}

func CompressString(input string) ([]byte, error) {
	if zkcommon.IsEmpty(input) {
		return nil, errors.New("empty input")
	}
	var b bytes.Buffer

	gzipWriter := gzip.NewWriter(&b)

	_, err := gzipWriter.Write([]byte(input))
	if err != nil {
		return nil, err
	}

	err = gzipWriter.Close()
	if err != nil {
		return nil, err
	}

	log.Println("String compressed successfully")

	return b.Bytes(), nil
}

func DecompressString(input []byte) (string, error) {
	if input == nil || len(input) == 0 {
		return "", errors.New("empty input")
	}

	reader := bytes.NewReader(input)

	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return "", err
	}
	defer gzipReader.Close()

	decompressed, err := io.ReadAll(gzipReader)
	if err != nil {
		return "", err
	}

	log.Println("String decompressed successfully")

	return string(decompressed), nil
}
