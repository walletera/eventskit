//go:build eventstoredb_test

package eventstoredb

import (
    "context"
    "fmt"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/walletera/eventskit/eventsourcing"
    "github.com/walletera/eventskit/messages"
    "github.com/walletera/eventskit/mocks/github.com/walletera/eventskit/events"
)

func TestNackRetries(t *testing.T) {
    ctx, _ := context.WithTimeout(context.Background(), testTimeout)

    client, err := GetESDBClient(eventStoreDBUrl)
    require.NoError(t, err)
    db := NewDB(client)
    eventDataMock := &events.MockEventData{}
    rawEvent := []byte(`{"thisIsA":"rawTestEvent"}`)
    eventDataMock.On("Serialize").Return(rawEvent, nil)
    eventDataMock.On("Type").Return("TestEventType")

    err = CreatePersistentSubscription(eventStoreDBUrl, "testStream", "testGroup")
    require.NoError(t, err)

    consumer, err := NewMessagesConsumer(eventStoreDBUrl, "testStream", "testGroup")
    require.NoError(t, err)
    messagesCh, err := consumer.Consume()
    require.NoError(t, err)

    // event can be appended
    werr := db.AppendEvents(ctx, "testStream", eventsourcing.ExpectedAggregateVersion{IsNew: true}, eventDataMock)
    require.NoError(t, werr)

    msg, err := waitForMessageWithTimeout(t, messagesCh, 2*time.Second)
    require.NoError(t, err)

    err = msg.Acknowledger().Nack(messages.NackOpts{
        Requeue:      true,
        MaxRetries:   1,
        ErrorCode:    0,
        ErrorMessage: "",
    })
    require.NoError(t, err)

    msgRetried, err := waitForMessageWithTimeout(t, messagesCh, 2*time.Second)
    require.NoError(t, err, "message was not retried")
    assert.Equal(t, rawEvent, msgRetried.Payload())

    err = msg.Acknowledger().Nack(messages.NackOpts{
        Requeue:      true,
        MaxRetries:   1,
        ErrorCode:    0,
        ErrorMessage: "",
    })
    require.NoError(t, err)

    // Message was parked
    retrievedEvents, err := db.ReadEvents(ctx, "$persistentsubscription-testGroup::testStream-parked")
    require.NoError(t, err)
    require.Len(t, retrievedEvents, 1, "the message was not parked")
    require.Equal(t, rawEvent, retrievedEvents[0].RawEvent, "the parked message is not the one we expected")
}

func waitForMessageWithTimeout(t *testing.T, messagesCh <-chan messages.Message, duration time.Duration) (messages.Message, error) {
    timeout := time.After(duration)
    select {
    case receivedMessage, ok := <-messagesCh:
        if !ok {
            t.Errorf("channel closed")
            return messages.Message{}, nil
        }
        return receivedMessage, nil
    case <-timeout:
        return messages.Message{}, fmt.Errorf("timeout waiting message")
    }
    return messages.Message{}, nil
}
