package badger

import (
	"bytes"
	"context"
	"errors"
	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/pb"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"time"
)

const (
	LogTag = "BadgerStoreHandler"
)

type (
	// DB defines an embedded key/value store database interface.
	BadgerDB interface {
		Get(key string) (value string, err error)
		Set(key, value string) error
		Has(key string) (bool, error)
		BulkSet(keyVals map[string]string, ttl int64) error
		BulkGetForPrefix(keyPrefix []string) (map[string]string, error)
	}
)

func (b *BadgerStoreHandler) Get(key string) (value string, err error) {
	var badgerDbValue []byte
	err = b.db.View(func(txn *badger.Txn) error {
		key := []byte(key)
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		// Accessing the value
		var value []byte
		err = item.Value(func(val []byte) error {
			value = append([]byte{}, val...)
			return nil
		})

		if err != nil {
			return err
		}

		badgerDbValue = make([]byte, len(value))
		copy(badgerDbValue, value) // Assign the value to temp
		zkLogger.ErrorF(LogTag, "Key: %s, Value: %s\n", key, string(badgerDbValue))
		return nil
	})

	if err != nil {
		return "", err
	}

	return string(badgerDbValue), nil
}

// Set implements the DB interface. It attempts to store a value for a given key
// and namespace. If the key/value pair cannot be saved, an error is returned.
func (b *BadgerStoreHandler) Set(key string, value string, ttl int64) error {
	err := b.db.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry([]byte(key), []byte(value)).WithTTL(time.Duration(ttl) * time.Second)
		err := txn.SetEntry(entry)
		if err != nil {
			zkLogger.ErrorF("Error in txn.SetEntry %s in BadgerDB: %s", key, err)
		}
		return err
	})
	if err != nil {
		zkLogger.ErrorF(LogTag, "Error while setting key %s in BadgerDB: %s", key, err)
		return err
	}
	return nil
}

// Has implements the DB interface. It returns a boolean reflecting if the
// datbase has a given key for a namespace or not. An error is only returned if
// an error to Get would be returned that is not of type badger.ErrKeyNotFound.
func (b *BadgerStoreHandler) Has(key string) (ok bool, err error) {
	_, err = b.Get(key)
	switch {
	case errors.Is(err, badger.ErrKeyNotFound):
		return false, nil
	case err == nil:
		return true, nil
	}
	zkLogger.ErrorF(LogTag, "Error while checking if key %s exists in BadgerDB: %s", key, err)
	return false, err
}

func (b *BadgerStoreHandler) BulkSet(keyVals map[string]string, ttl int64) error {
	bulkWriter := b.db.NewWriteBatch()

	for key, val := range keyVals {
		if err := bulkWriter.SetEntry(badger.NewEntry([]byte(key), []byte(val)).
			WithTTL(time.Duration(ttl) * time.Second)); err != nil {
			return err
		}
	}

	err := bulkWriter.Flush()
	if err != nil {
		zkLogger.Error(badgerHandlerLogTag, "Error while syncing data to Badger ", err)
		return err
	}

	return nil
}

func (b *BadgerStoreHandler) BulkGetForPrefix(keyPrefix []string) (map[string]string, error) {

	resultSet := make(map[string]string)
	stream := b.db.NewStream()
	// db.NewStreamAt(readTs) for managed mode.

	// -- Optional settings
	stream.NumGo = 16 // Set number of goroutines to use for iteration.

	// ChooseKey is called concurrently for every key. If left nil, assumes true by default.
	stream.ChooseKey = func(item *badger.Item) bool {
		for _, key := range keyPrefix {
			if bytes.HasPrefix(item.Key(), []byte(key)) {
				return true
			}
		}
		return false
	}
	stream.KeyToList = nil

	// Send is called serially, while Stream.Orchestrate is running.
	stream.Send = func(list *pb.KVList) error {
		recordItems := list.GetKv()
		for _, record := range recordItems {
			resultSet[string(record.GetKey())] = string(record.GetValue())
		}
		return nil
	}

	// Run the stream
	if err := stream.Orchestrate(context.Background()); err != nil {
		return resultSet, err
	}

	return resultSet, nil
}
