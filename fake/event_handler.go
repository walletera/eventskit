package fake

import (
    "context"

    "github.com/walletera/eventskit/errors"
)

type EventHandler interface {
    HandleFakeEvent(ctx context.Context, e Event) errors.ProcessingError
}
