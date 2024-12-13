package payments

import "github.com/walletera/eventskit/messages"

// NewProcessor returns a messages.Processor which is specific to payments events
func NewProcessor(messageConsumer messages.Consumer, eventsHandler EventsHandler, opts ...messages.ProcessorOpt) *messages.Processor[EventsHandler] {
    return messages.NewProcessor[EventsHandler](
        messageConsumer,
        NewEventsDeserializer(),
        eventsHandler,
        opts...,
    )
}
