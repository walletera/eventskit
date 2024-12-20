package fake

import (
    "context"
    "encoding/json"

    "github.com/walletera/eventskit/events"
    "github.com/walletera/werrors"
)

var _ events.Event[EventHandler] = Event{}

type Event struct {
    FakeID              string `json:"id"`
    FakeType            string `json:"type"`
    FakeCorrelationID   string `json:"correlation_id"`
    FakeDataContentType string `json:"data_content_type"`
    FakeData            string `json:"data"`
}

func (f Event) ID() string {
    return f.FakeID
}

func (f Event) Type() string {
    return f.FakeType
}

func (f Event) CorrelationID() string {
    return f.FakeCorrelationID
}

func (f Event) DataContentType() string {
    return f.FakeDataContentType
}

func (f Event) Serialize() ([]byte, error) {
    serialized, err := json.Marshal(f)
    if err != nil {
        return nil, err
    }
    return serialized, nil
}

func (f Event) Accept(ctx context.Context, visitor EventHandler) werrors.WError {
    return visitor.HandleFakeEvent(ctx, f)
}
