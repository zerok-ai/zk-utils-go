package zktick

import (
	"fmt"
	"time"
)

// TickerTask create a struct to hold the ticker and the task
type TickerTask struct {
	ticker   *time.Ticker
	name     string
	task     func()
	counter  int
	interval time.Duration
}

func GetNewTickerTask(name string, interval time.Duration, task func()) *TickerTask {
	return &TickerTask{
		ticker:   time.NewTicker(interval),
		task:     task,
		name:     name,
		interval: interval,
	}
}

func (tt TickerTask) Start() *TickerTask {
	go func() {
		runTaskAndFlushPendingTasks(tt)
		time.AfterFunc(tt.interval, func() {
			for {
				select {
				case <-tt.ticker.C:
					// Perform the task
					fmt.Printf("tick (%s) - %d\n", tt.name, tt.counter)
					tt.counter = tt.counter + 1
					runTaskAndFlushPendingTasks(tt)
				}
			}
		})
	}()
	return &tt
}

func runTaskAndFlushPendingTasks(tt TickerTask) {
	tt.task()
	//If there are multiple ticks available, flush all for now
	for len(tt.ticker.C) > 0 {
		<-tt.ticker.C
		fmt.Println("Skipping tick due to slow processing")
	}
}

func (tt TickerTask) Stop() *TickerTask {
	tt.ticker.Stop()
	return &tt
}
