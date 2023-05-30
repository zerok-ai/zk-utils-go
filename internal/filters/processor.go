package filters

import (
	"time"

	"scenario-manager/internal/config"
	zktime "scenario-manager/zk-utils-go/time"
)

const (
	filterProcessingTickInterval time.Duration = 10 * time.Second
)

var (
	filterProcessor *FilterProcessor

	//	tickers
	tickerTraceProcessor *time.Ticker
)

func Start(cfg config.AppConfigs) error {

	// initialize the filter store
	var err error
	if filterProcessor, err = NewFilterProcessor(cfg); err != nil {
		return err
	}

	// trigger recurring processing of trace data against filters
	tickerTraceProcessor = time.NewTicker(filterProcessingTickInterval)
	zktime.RunTaskOnTicks(tickerTraceProcessor, processFilters)
	processFilters()

	populateData()

	return nil
}

func processFilters() {

}
