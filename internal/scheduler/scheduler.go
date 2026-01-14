package scheduler

import (
	"context"
	"log"
	"strconv"
	"time"

	"tg-bots/internal/storage"
)

// UserLister provides access to known user IDs.
type UserLister interface {
	ListUserIDs() []int64
}

// Sender sends a reminder message to a user.
type Sender func(userID int64, text string) error

// Run starts the reminder loop until the context is canceled.
func Run(ctx context.Context, interval time.Duration, store storage.Store, users UserLister, send Sender) {
	if interval <= 0 {
		interval = 20 * time.Second
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	runOnce := func() {
		userIDs := users.ListUserIDs()
		now := time.Now().UTC()
		for _, userID := range userIDs {
			tasks, err := store.ListTasks(userID)
			if err != nil {
				log.Printf("list tasks failed user=%d: %v", userID, err)
				continue
			}
			for _, task := range tasks {
				if task.Done || task.ReminderAt == nil {
					continue
				}
				if task.ReminderAt.After(now) {
					continue
				}
				text := "Напоминание (#" + itoa(task.ID) + "): " + task.Text
				if err := send(userID, text); err != nil {
					log.Printf("send reminder failed user=%d task=%d: %v", userID, task.ID, err)
					continue
				}
				// Clear reminder so it fires once.
				task.ReminderAt = nil
				if err := store.UpdateTask(userID, task); err != nil {
					log.Printf("update task after reminder failed user=%d task=%d: %v", userID, task.ID, err)
				}
			}
		}
	}

	runOnce()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			runOnce()
		}
	}
}

func itoa(v int) string {
	return strconv.Itoa(v)
}
