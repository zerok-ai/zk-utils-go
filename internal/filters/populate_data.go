package filters

import (
	zktime "github.com/zerok-ai/zk-utils-go/ticker"
	"time"
)

const runseq = "R1"

var (
	counter = 0
)

func populateData() {
	tickerTraceProcessor = time.NewTicker(1 * time.Second)
	zktime.RunTaskOnTicks(tickerTraceProcessor, populateDataOddKeys)

	tickerTraceProcessor = time.NewTicker(2 * time.Second)
	zktime.RunTaskOnTicks(tickerTraceProcessor, populateDataOddKeys)

}

func populateDataOddKeys() {
	//store := filterProcessor.versionedStore
	//store.SetValue(fmt.Sprintf("%s%n", runseq, counter), fmt.Sprintf("%s%n", runseq, counter))
}
