//go:build eventstoredb_test

package eventstoredb

import (
    "context"
    "testing"

    "github.com/EventStore/EventStore-Client-Go/v4/esdb"
    "github.com/stretchr/testify/require"
    "github.com/walletera/eventskit/eventsourcing"
    "github.com/walletera/eventskit/mocks/github.com/walletera/eventskit/events"
    "github.com/walletera/werrors"
)

func TestAppendReadEvents(t *testing.T) {
    const streamName = "testStream"

    client, err := GetESDBClient("esdb://localhost:2113?tls=false")
    require.NoError(t, err)

    ctx, _ := context.WithTimeout(context.Background(), testTimeout)

    t.Cleanup(func() {
        _, deleteStreamErr := client.DeleteStream(ctx, streamName, esdb.DeleteStreamOptions{})
        require.NoError(t, deleteStreamErr)
    })

    db := NewDB(client)
    eventDataMock := &events.MockEventData{}
    rawEvent := []byte(`{"thisIsA":"rawTestEvent"}`)
    eventDataMock.On("Serialize").Return(rawEvent, nil)
    eventDataMock.On("Type").Return("TestEventType")

    // event can be appended
    nextExpectedVersion, werr := db.AppendEvents(ctx, streamName, eventsourcing.ExpectedAggregateVersion{IsNew: true}, eventDataMock)
    require.NoError(t, werr)
    require.Equal(t, uint64(0), nextExpectedVersion)

    // try to append again specifying eventsourcing.ExpectedAggregateVersion{IsNew: true}
    // should result in a ResourceAlreadyExist error
    _, werr = db.AppendEvents(ctx, streamName, eventsourcing.ExpectedAggregateVersion{IsNew: true}, eventDataMock)
    require.Equal(t, werrors.ResourceAlreadyExistErrorCode, werr.Code())

    // event can be read
    retrievedEvents, err := db.ReadEvents(ctx, streamName)
    require.NoError(t, err)
    require.Len(t, retrievedEvents, 1)
    require.Equal(t, rawEvent, retrievedEvents[0].RawEvent)
    require.Equal(t, uint64(0), retrievedEvents[0].AggregateVersion)

    // a new event can be appended to the stream with correct eventsourcing.ExpectedAggregateVersion
    nextExpectedVersion, werr = db.AppendEvents(ctx, streamName, eventsourcing.ExpectedAggregateVersion{Version: 0}, eventDataMock)
    require.NoError(t, werr)

    // the next expected aggregate version is previous versions + 1 because we appended 1 event
    require.Equal(t, uint64(1), nextExpectedVersion)

    // try to append a new event with a wrong version result in a werrors.WrongAggregateVersionErrorCode error
    _, werr = db.AppendEvents(ctx, streamName, eventsourcing.ExpectedAggregateVersion{Version: 0}, eventDataMock)
    require.Equal(t, werrors.WrongResourceVersionErrorCode, werr.Code())
}
