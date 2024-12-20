package tests

import (
    "context"
    "fmt"
    "sync"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
    "github.com/walletera/eventskit/fake"
    "github.com/walletera/eventskit/messages"
    eventsmock "github.com/walletera/eventskit/mocks/github.com/walletera/eventskit/events"
    fakemock "github.com/walletera/eventskit/mocks/github.com/walletera/eventskit/fake"
    messagesmock "github.com/walletera/eventskit/mocks/github.com/walletera/eventskit/messages"
    "github.com/walletera/werrors"
)

func TestMessageProcessor(t *testing.T) {

    var testErrs = map[string]werrors.WError{
        "noErr":                     nil,
        "unprocessableMessageError": werrors.NewUnprocessableMessageError("boom"),
        "internalError":             werrors.NewNonRetryableInternalError("boom"),
        "timeoutError":              werrors.NewTimeoutError("boom"),
    }

    for errName, err := range testErrs {
        var acknowledgerMockExpectationSetter func(mock *messagesmock.MockAcknowledger)
        var testName string
        if errName == "noErr" {
            acknowledgerMockExpectationSetter = func(mock *messagesmock.MockAcknowledger) {
                mock.On("Ack").Return(nil)
            }
            testName = "if no error then the message is Ack"
        } else {
            acknowledgerMockExpectationSetter = func(mock *messagesmock.MockAcknowledger) {
                mock.On("Nack", messages.NackOpts{
                    Requeue:      err.IsRetryable(),
                    ErrorCode:    err.Code(),
                    ErrorMessage: err.Message(),
                }).Return(nil)
            }
            testName = fmt.Sprintf("if %s occurred then message is Nack", errName)
        }

        t.Run(testName, func(t *testing.T) {
            execTest(t, acknowledgerMockExpectationSetter, err)
        })
    }
}

func execTest(t *testing.T, acknowledgerMockExpectationSetter func(acknowledgerMock *messagesmock.MockAcknowledger), err werrors.WError) {
    messagesCh := make(chan messages.Message)

    messageConsumerMock := &messagesmock.MockConsumer{}
    messageConsumerMock.On("Consume").Return((<-chan messages.Message)(messagesCh), nil)
    messageConsumerMock.On("Close").Return(nil).Run(func(_ mock.Arguments) {
        close(messagesCh)
    })

    acknowledgerMock := &messagesmock.MockAcknowledger{}
    acknowledgerMockExpectationSetter(acknowledgerMock)

    rawPayload := []byte("raw message payload")
    message := messages.NewMessage(rawPayload, acknowledgerMock)

    event := fake.Event{}

    eventsDeserializerMock := &eventsmock.MockDeserializer[fake.EventHandler]{}
    eventsDeserializerMock.On("Deserialize", rawPayload).Return(event, nil)

    wg := sync.WaitGroup{}
    wg.Add(1)
    mockFakeEventHandler := &fakemock.MockEventHandler{}
    mockFakeEventHandler.
        On("HandleFakeEvent", mock.Anything, event).
        Return(err).
        Run(func(args mock.Arguments) {
            wg.Done()
        })

    messageProcessor := messages.NewProcessor[fake.EventHandler](
        messageConsumerMock,
        eventsDeserializerMock,
        mockFakeEventHandler,
    )

    ctx, cancel := context.WithCancel(context.Background())

    messageProcessorStartError := messageProcessor.Start(ctx)
    require.NoError(t, messageProcessorStartError)

    messagesCh <- message

    wg.Wait()

    cancel()

    _, open := <-messagesCh
    assert.False(t, open, "the messages channel is still open")

    acknowledgerMock.AssertExpectations(t)
    messageConsumerMock.AssertExpectations(t)
    eventsDeserializerMock.AssertExpectations(t)
    mockFakeEventHandler.AssertExpectations(t)
}
