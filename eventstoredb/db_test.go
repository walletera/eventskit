//go:build eventstoredb_test

package eventstoredb

import (
    "context"
    "fmt"
    "net/http"
    "testing"
    "time"

    "github.com/stretchr/testify/require"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
    "github.com/walletera/eventskit/containertest"
    "github.com/walletera/eventskit/eventsourcing"
    "github.com/walletera/eventskit/mocks/github.com/walletera/eventskit/events"
    "github.com/walletera/werrors"
)

const (
    eventStoreDBPort            = "2113"
    testTimeout                 = 30 * time.Second
    containerStartTimeout       = 10 * time.Second
    containerTerminationTimeout = 10 * time.Second
    eventStoreDBUrl             = "esdb://localhost:2113?tls=false"
)

func TestAppendReadEvents(t *testing.T) {
    ctx, _ := context.WithTimeout(context.Background(), testTimeout)
    terminateContainer, err := startEventStoreDBContainer(ctx)
    defer terminateContainer()
    require.NoError(t, err)
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
    require.Equal(t, werrors.AggregateAlreadyExistErrorCode, werr.Code())
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
    require.Equal(t, werrors.WrongAggregateVersionErrorCode, werr.Code())
}

func startEventStoreDBContainer(ctx context.Context) (func() error, error) {
    req := testcontainers.ContainerRequest{
        Image: "eventstore/eventstore:21.10.7-buster-slim",
        Name:  "esdb-node",
        Cmd:   []string{"--insecure", "--run-projections=All"},
        ExposedPorts: []string{
            fmt.Sprintf("%s:%s", eventStoreDBPort, eventStoreDBPort),
        },
        WaitingFor: wait.
            ForHTTP("/health/live").
            WithPort("2113/tcp").
            WithStartupTimeout(containerStartTimeout).
            WithStatusCodeMatcher(func(status int) bool {
                return status == http.StatusNoContent
            }),
        LogConsumerCfg: &testcontainers.LogConsumerConfig{
            Consumers: []testcontainers.LogConsumer{containertest.NewContainerLogConsumer("esdb")},
        },
    }
    esdbContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    if err != nil {
        return nil, fmt.Errorf("error creating esdb container: %w", err)
    }

    return func() error {
        terminationCtx, terminationCtxCancel := context.WithTimeout(context.Background(), containerTerminationTimeout)
        defer terminationCtxCancel()
        terminationErr := esdbContainer.Terminate(terminationCtx)
        if terminationErr != nil {
            fmt.Errorf("failed terminating esdb container: %w", err)
        }
        return nil
    }, nil
}
