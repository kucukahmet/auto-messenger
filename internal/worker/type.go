package worker

import (
	"auto-messager/internal/app"
	"sync"
)

type Listener struct {
	isRunning bool
	app       *app.App
	mutex     sync.Mutex
	waitGroup sync.WaitGroup
	stopChan  chan struct{}
}
