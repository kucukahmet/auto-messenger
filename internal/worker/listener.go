package worker

import (
	"auto-messager/internal/app"
	"auto-messager/internal/storage"
	"auto-messager/internal/utils"
	"context"
	"fmt"
	"log"
	"sync"
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	limit := int32(listener.app.Config.BATCH_SIZE)
	messages, err := listener.app.Queries.GetPendingForUpdate(ctx, limit)

	if err != nil {
		fmt.Printf("Error fetching pending messages: %v\n", err)
		return
	}

	if len(messages) == 0 {
		log.Println("No pending messages to process")
		return
	}

	fmt.Printf("Fetched %d pending messages\n", len(messages))

	var waitGroup sync.WaitGroup
	for _, message := range messages {

		if listener.app.Cache.Get(ctx, fmt.Sprintf("message:%d", message.ID)).Val() != "" {
			fmt.Printf("Message ID %d is already being processed, skipping...\n", message.ID)
			continue
		}

		listener.app.Cache.Set(ctx, fmt.Sprintf("message:%d", message.ID), "processing", 10*time.Minute)
		waitGroup.Add(1)
		go func(msg storage.Message) {
			defer waitGroup.Done()

			payload, err := utils.BuildPayloadFromMessage(&msg)
			if err != nil {
				fmt.Printf("Error building payload for message ID %d: %v\n", msg.ID, err)
				return
			}

			response, err := listener.app.Service.SendMessage(payload)
			if err != nil {
				fmt.Printf("Error sending message ID %d: %v\n", msg.ID, err)
				return
			}

			wresponse, err := utils.ParseWebhookResponse(response)
			if err != nil {
				fmt.Printf("Error parsing response for message ID %d: %v\n", msg.ID, err)
				return
			}
			listener.app.Queries.MarkSent(ctx, storage.MarkSentParams{
				ID:                msg.ID,
				ResponseMessageID: wresponse.MessageID,
			})
			listener.app.Cache.Set(ctx, fmt.Sprintf("message:%d", msg.ID), wresponse.MessageID, 10*time.Minute)
			fmt.Printf("Sent message ID %d, response status: %v -> %v\n", msg.ID, wresponse.Message, wresponse.MessageID)

			fmt.Println("Processing message ID:", msg.ID)
		}(message)
	}
	waitGroup.Wait()
	fmt.Println("All messages processed")
}
