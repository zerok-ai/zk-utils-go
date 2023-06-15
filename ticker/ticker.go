package zktick

import (
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"time"
)

var LogTag = "zkTick_ticker"

// TickerTask create a struct to hold the ticker and the task
type TickerTask struct {
	ticker  *time.Ticker
	name    string
	task    func()
	counter int
}

func GetNewTickerTask(name string, interval time.Duration, task func()) *TickerTask {
	return &TickerTask{
		ticker: time.NewTicker(interval),
		task:   task,
		name:   name,
	}
}

func (tt TickerTask) Start() {
	go func() {
		for {
			select {
			case <-tt.ticker.C:
				// Perform the task
				zkLogger.Error(LogTag, "tick (%s) - %d\n", tt.name, tt.counter)
				tt.counter = tt.counter + 1
				tt.task()

				//If there are multiple ticks available, flush all for now
				for len(tt.ticker.C) > 0 {
					<-tt.ticker.C
					zkLogger.Info(LogTag, "Skipping tick due to slow processing")
				}
			}
		}
	}()
}

func (tt TickerTask) Stop() {
	tt.ticker.Stop()
}
