package payments

import (
    "context"

    "github.com/walletera/eventskit/errors"
)

type EventsHandler interface {
    HandleWithdrawalCreated(ctx context.Context, withdrawalCreated WithdrawalCreatedEvent) errors.ProcessingError
}
