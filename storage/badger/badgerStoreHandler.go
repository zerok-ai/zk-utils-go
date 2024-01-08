package badger

import (
	"context"
	"github.com/kataras/iris/v12/x/errors"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/storage/badger/config"
	zktick "github.com/zerok-ai/zk-utils-go/ticker"
	"time"

	"github.com/dgraph-io/badger"
)

const (
	// Default BadgerDB discardRatio. It represents the discard ratio for the
	// BadgerDB GC.
	//
	// Ref: https://godoc.org/github.com/dgraph-io/badger#DB.RunValueLogGC
	badgerDiscardRatio = 0.5

	// Default BadgerDB GC interval
	badgerGCInterval = 10 * time.Minute

	// Log tag
	badgerHandlerLogTag = "BadgerStoreHandler"
)

var (
	// BadgerAlertNamespace defines the alerts BadgerDB namespace.
	BadgerAlertNamespace = []byte("alerts")
)

type (
	// BadgerStoreHandler is a wrapper around a BadgerDB backend database that implements
	// the DB interface.
	BadgerStoreHandler struct {
		ctx          context.Context
		badgerConfig *config.BadgerConfig
		gcTicker     *zktick.TickerTask
		db           *badger.DB
	}
)

// NewBadgerDB returns a new initialized BadgerDB database implementing the DB
// interface. If the database cannot be initialized, an error will be returned.
func NewBadgerHandler(configs *config.BadgerConfig) (*BadgerStoreHandler, error) {
	handler := BadgerStoreHandler{
		ctx:          context.Background(),
		badgerConfig: configs,
	}

	err := handler.InitializeConn()
	if err != nil {
		zkLogger.Error(badgerHandlerLogTag, "Error while initializing connection ", err)
		return nil, err
	}

	timerDuration := time.Duration(configs.GCTimerDuration) * time.Second
	handler.gcTicker = zktick.GetNewTickerTask("badger_garbage_collect", timerDuration, handler.runGC)
	handler.gcTicker.Start()
	handler.ctx = context.Background()

	return &handler, nil
}

func (b *BadgerStoreHandler) InitializeConn() error {

	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	db, err := badger.Open(badger.DefaultOptions(b.badgerConfig.DBPath))
	if err != nil {
		zkLogger.Error(badgerHandlerLogTag, "Error while initializing connection ", err)
		return err
	}
	b.db = db
	return nil
}

// Close implements the DB interface. It closes the connection to the underlying
// BadgerDB database as well as invoking the context's cancel function.
func (b *BadgerStoreHandler) Close() error {
	b.gcTicker.Stop()
	return b.db.Close()
}

// runGC triggers the garbage collection for the BadgerDB backend database. It
// should be run in a goroutine.
func (b *BadgerStoreHandler) runGC() {
	b.StartCompaction()
	err := b.db.RunValueLogGC(b.badgerConfig.GCDiscardRatio)
	if err != nil {
		if errors.Is(err, badger.ErrNoRewrite) {
			zkLogger.Debug(badgerHandlerLogTag, "No BadgerDB GC occurred:", err)
		} else {
			zkLogger.Error(badgerHandlerLogTag, "Error while running garbage collector ", err)
		}
	}
}

func (b *BadgerStoreHandler) StartCompaction() {
	err := b.db.Flatten(2)
	if err != nil {
		if errors.Is(err, badger.ErrNoRewrite) {
			zkLogger.Debug(badgerHandlerLogTag, "No BadgerDB GC occurred:", err)
		} else {
			zkLogger.Error(badgerHandlerLogTag, "Error while running garbage collector ", err)
		}
	}
}
