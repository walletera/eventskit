//go:build eventstoredb_test

package eventstoredb

import (
    "context"
    "testing"

    "github.com/stretchr/testify/require"
    "github.com/walletera/eventskit/eventsourcing"
    "github.com/walletera/eventskit/mocks/github.com/walletera/eventskit/events"
    "github.com/walletera/werrors"
)

func TestAppendReadEvents(t *testing.T) {
    ctx, _ := context.WithTimeout(context.Background(), testTimeout)

    client, err := GetESDBClient(eventStoreDBUrl)
    require.NoError(t, err)
    db := NewDB(client)
    eventDataMock := &events.MockEventData{}
    rawEvent := []byte(`{"thisIsA":"rawTestEvent"}`)
    eventDataMock.On("Serialize").Return(rawEvent, nil)
    eventDataMock.On("Type").Return("TestEventType")

    // event can be appended
    werr := db.AppendEvents(ctx, "testStream", eventsourcing.ExpectedAggregateVersion{IsNew: true}, eventDataMock)
    require.NoError(t, werr)

    // try to append again specifying eventsourcing.ExpectedAggregateVersion{IsNew: true}
    // should result in a ResourceAlreadyExist error
    werr = db.AppendEvents(ctx, "testStream", eventsourcing.ExpectedAggregateVersion{IsNew: true}, eventDataMock)
    require.Equal(t, werrors.ResourceAlreadyExistErrorCode, werr.Code())

    // event can be read
    retrievedEvents, err := db.ReadEvents(ctx, "testStream")
    require.NoError(t, err)
    require.Len(t, retrievedEvents, 1)
    require.Equal(t, rawEvent, retrievedEvents[0].RawEvent)
    require.Equal(t, uint64(0), retrievedEvents[0].AggregateVersion)

    // new event can be appended to the stream with correct eventsourcing.ExpectedAggregateVersion
    werr = db.AppendEvents(ctx, "testStream", eventsourcing.ExpectedAggregateVersion{Version: 0}, eventDataMock)
    require.NoError(t, werr)

    // try to append new event with wrong version result in a werrors.WrongAggregateVersionErrorCode error
    werr = db.AppendEvents(ctx, "testStream", eventsourcing.ExpectedAggregateVersion{Version: 0}, eventDataMock)
    require.Equal(t, werrors.WrongResourceVersionErrorCode, werr.Code())
}
