package eventsourcing

import (
    "context"

    "github.com/walletera/eventskit/events"
)

type RawEvent []byte

type DB interface {
    AppendEvents(ctx context.Context, streamName string, event ...events.EventData) error
    ReadEvents(ctx context.Context, streamName string) ([]RawEvent, error)
}
