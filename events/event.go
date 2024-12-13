package events

import (
    "context"

    "github.com/walletera/eventskit/errors"
)

type Event[Handler any] interface {
    EventData

    Accept(ctx context.Context, Handler Handler) errors.ProcessingError
}
