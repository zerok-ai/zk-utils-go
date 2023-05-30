package filters

import (
	"github.com/zerok-ai/zk-utils-go/rules/model"
	"github.com/zerok-ai/zk-utils-go/storage"
)

type StoreType model.Scenario

type FilterProcessor struct {
	versionedStore *storage.VersionedStore[StoreType]
}
