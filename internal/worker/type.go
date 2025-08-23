package worker

import "sync"

type Listener struct {
	isRunning bool
	period    int32
	mutex     sync.Mutex
	waitGroup sync.WaitGroup
	stopChan  chan struct{}
}
