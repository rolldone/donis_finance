package donisfinance

import (
	"context"
	"log"

	"github.com/rolldone/donisgo/internal/events"
)

// registerEventHandlers registers example event handlers for the plugin.
func registerEventHandlers() {
	// subscribe to a sample event; handlers run asynchronously
		events.Subscribe("user.created", func(ctx context.Context, payload interface{}) {
			log.Printf("plugin donisfinance: received user.created payload type=%T", payload)
			_ = ctx
		})
}
