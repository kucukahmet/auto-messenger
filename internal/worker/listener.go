package worker

import (
	"fmt"
	"time"
)

func NewListener(period int32) *Listener {
	return &Listener{
		isRunning: false,
		period:    period,
		stopChan:  make(chan struct{}),
	}
}

func (listener *Listener) IsRunning() bool {
	listener.mutex.Lock()
	defer listener.mutex.Unlock()
	return listener.isRunning
}

func (listener *Listener) Start() {
	listener.mutex.Lock()
	defer listener.mutex.Unlock()
	if listener.isRunning {
		return
	}
	listener.isRunning = true
	listener.stopChan = make(chan struct{})
	listener.waitGroup.Add(1)

	go func() {
		defer listener.waitGroup.Done()
		ticker := time.NewTicker(time.Duration(listener.period) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fmt.Println("Listener tick ...")

			case <-listener.stopChan:
				fmt.Println("Listener stopping ...")
				return
			}
		}
	}()
}

func (listener *Listener) Stop() {
	listener.mutex.Lock()
	defer listener.mutex.Unlock()
	if listener.isRunning {
		return
	}
	listener.isRunning = false
	close(listener.stopChan)
}
