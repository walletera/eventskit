package main

import (
    "context"
    "log"

    "github.com/walletera/eventskit/errors"
    "github.com/walletera/eventskit/messages"
    "github.com/walletera/eventskit/payments"
)

func main() {

    paymentsEventsHandler := NewPaymentsEventsHandlerImpl()

    processor, err := payments.NewRabbitMQProcessor(
        paymentsEventsHandler,
        "message-processor-example-queue",
        payments.RabbitMQProcessorOpt{
            ProcessorOpt: messages.WithErrorCallback(
                func(processingError errors.ProcessingError) {
                    log.Printf("error processing message: %s", processingError.Error())
                },
            ),
        },
    )
    if err != nil {
        log.Fatalf("error creating new RabbitMQProcessor: %s", err.Error())
    }

    err = processor.Start(context.Background())
    if err != nil {
        log.Fatalf("failed to start message processor: %s", err.Error())
    }

    blockForeverCh := make(chan any)
    <-blockForeverCh
}

type PaymentsEventsHandlerImpl struct {
}

func NewPaymentsEventsHandlerImpl() *PaymentsEventsHandlerImpl {
    return &PaymentsEventsHandlerImpl{}
}

func (p PaymentsEventsHandlerImpl) HandleWithdrawalCreated(_ context.Context, withdrawalCreated payments.WithdrawalCreatedEvent) errors.ProcessingError {
    log.Printf("handling WithdrawalCreatedEvent event: %+v", withdrawalCreated)
    return nil
}
