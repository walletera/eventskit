package fake

import (
    "context"

    "github.com/walletera/werrors"
)

type EventHandler interface {
    HandleFakeEvent(ctx context.Context, e Event) werrors.WError
}
