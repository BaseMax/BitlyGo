package schedulers

import (
	"context"
	"fmt"
	"time"

	"github.com/itsjoniur/bitlygo/internal/durable"
)

func RemoveExpiredLinks(db *durable.Database, logger *durable.Logger) {
	ticker := time.NewTicker(time.Minute * 1)
	done := make(chan bool)

	go func() {
		query := "DELETE FROM links WHERE expired_at < NOW() - INTERVAL '1 minute';"
		_, err := db.Exec(context.Background(), query)
		if err != nil {
			logger.Error(fmt.Sprintf("Can not remove expired links: %s", err.Error()))
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			continue
		}
	}
}
