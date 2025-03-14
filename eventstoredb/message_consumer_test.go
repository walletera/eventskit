//go:build eventstoredb_test

package eventstoredb

import (
    "context"
    "fmt"
    "testing"
    "time"

    "github.com/EventStore/EventStore-Client-Go/v4/esdb"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/walletera/eventskit/eventsourcing"
    "github.com/walletera/eventskit/messages"
    "github.com/walletera/eventskit/mocks/github.com/walletera/eventskit/events"
)

func TestNackRetries(t *testing.T) {
    ctx, _ := context.WithTimeout(context.Background(), testTimeout)

    client, err := GetESDBClient("esdb://localhost:2113?tls=false")
    require.NoError(t, err)
    db := NewDB(client)
    eventDataMock := &events.MockEventData{}
    rawEvent := []byte(`{"thisIsA":"rawTestEvent"}`)
    eventDataMock.On("Serialize").Return(rawEvent, nil)
    eventDataMock.On("Type").Return("TestEventType")

    maxRetries := 3
    subscriptionSettings := esdb.SubscriptionSettingsDefault()
    subscriptionSettings.ResolveLinkTos = true
    subscriptionSettings.MaxRetryCount = int32(maxRetries)

    err = CreatePersistentSubscription(eventStoreDBUrl, "testStream", "testGroup", subscriptionSettings)
    require.NoError(t, err)

    consumer, err := NewMessagesConsumer(eventStoreDBUrl, "testStream", "testGroup")
    require.NoError(t, err)
    messagesCh, err := consumer.Consume()
    require.NoError(t, err)

    info, err := db.client.GetPersistentSubscriptionInfo(ctx, "testStream", "testGroup", esdb.GetPersistentSubscriptionOptions{})
    require.NoError(t, err)

    // parked messages count is 0
    require.EqualValues(t, 0, info.Stats.ParkedMessagesCount)

    // event can be appended
    werr := db.AppendEvents(ctx, "testStream", eventsourcing.ExpectedAggregateVersion{IsNew: true}, eventDataMock)
    require.NoError(t, werr)

    msg, err := waitForMessageWithTimeout(t, messagesCh, 2*time.Second)
    require.NoError(t, err)

    nackOpts := messages.NackOpts{
        Requeue: true,
    }

    err = msg.Acknowledger().Nack(nackOpts)
    require.NoError(t, err)

    for i := 0; i < maxRetries; i++ {
        msgRetried, err := waitForMessageWithTimeout(t, messagesCh, 2*time.Second)
        require.NoError(t, err, "message was not retried %d times", i)
        assert.Equal(t, rawEvent, msgRetried.Payload())

        err = msgRetried.Acknowledger().Nack(nackOpts)
        require.NoError(t, err)
    }

    // let's give a moment to the stats to be updated
    time.Sleep(100 * time.Millisecond)

    info, err = db.client.GetPersistentSubscriptionInfo(ctx, "testStream", "testGroup", esdb.GetPersistentSubscriptionOptions{})
    require.NoError(t, err)
    require.EqualValues(t, 1, info.Stats.ParkedMessagesCount, "the message was not parked")

    // replay parked message
    err = db.client.ReplayParkedMessages(ctx, "testStream", "testGroup", esdb.ReplayParkedMessagesOptions{})
    require.NoError(t, err)

    msgReplayed, err := waitForMessageWithTimeout(t, messagesCh, 2*time.Second)
    require.NoError(t, err, "message was not replayed")
    assert.Equal(t, rawEvent, msgReplayed.Payload())

    // let's give a moment to the stats to be updated
    time.Sleep(100 * time.Millisecond)

    info, err = db.client.GetPersistentSubscriptionInfo(ctx, "testStream", "testGroup", esdb.GetPersistentSubscriptionOptions{})
    require.NoError(t, err)
    require.EqualValues(t, 0, info.Stats.ParkedMessagesCount, "the message is still parked")
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
