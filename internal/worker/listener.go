package worker

import (
	"auto-messager/internal/app"
	"fmt"
	"time"
)

func NewListener(app *app.App) *Listener {
	return &Listener{
		isRunning: false,
		app:       app,
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
		ticker := time.NewTicker(time.Duration(listener.app.Config.PERIOD) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fmt.Println("Listener tick ...")
				listener.action()

			case <-listener.stopChan:
				fmt.Println("Listener stopping ...")
				return
			}
		}
	}()
}

func (listener *Listener) Stop() error {
	listener.mutex.Lock()
	defer listener.mutex.Unlock()
	if !listener.isRunning {
		return fmt.Errorf("listener is not running")
	}
	close(listener.stopChan)
	listener.waitGroup.Wait()
	listener.isRunning = false
	return nil
}

func (listener *Listener) action() {
	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()

	// storage.GetPendingForUpdate(ctx, listener.period)

}
