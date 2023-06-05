package crypto

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"github.com/google/uuid"
	"io"
	"log"
)

func CalculateHash(s string) uuid.UUID {
	hash := sha1.Sum([]byte(s))
	return uuid.NewSHA1(uuid.Nil, hash[:])
}

func CompressString(input string) ([]byte, error) {
	var b bytes.Buffer

	// Create a gzip writer on top of the buffer
	gzipWriter := gzip.NewWriter(&b)

	// Write the input string to the gzip writer
	_, err := gzipWriter.Write([]byte(input))
	if err != nil {
		return nil, err
	}

	// Close the gzip writer to flush any remaining data
	err = gzipWriter.Close()
	if err != nil {
		return nil, err
	}

	log.Println("String compressed successfully")

	return b.Bytes(), nil
}

func DecompressString(input []byte) (string, error) {
	// Create a bytes reader from the input
	reader := bytes.NewReader(input)

	// Create a gzip reader on top of the bytes reader
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return "", err
	}
	defer gzipReader.Close()

	// Read the decompressed data from the gzip reader
	decompressed, err := io.ReadAll(gzipReader)
	if err != nil {
		return "", err
	}

	log.Println("String decompressed successfully")

	return string(decompressed), nil
}
