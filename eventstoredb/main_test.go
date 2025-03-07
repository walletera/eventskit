package eventstoredb

import (
    "context"
    "fmt"
    "net/http"
    "testing"
    "time"

    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
    "github.com/walletera/eventskit/containertest"
)

const (
    eventStoreDBPort            = "2113"
    testTimeout                 = 10 * time.Second
    containerStartTimeout       = 10 * time.Second
    containerTerminationTimeout = 10 * time.Second
    eventStoreDBUrl             = "esdb://localhost:2113?tls=false"
)

func TestMain(m *testing.M) {
    ctx, _ := context.WithTimeout(context.Background(), containerStartTimeout)
    terminateContainer, err := startEventStoreDBContainer(ctx)
    if err != nil {
        panic("error starting eventstoreDB container: " + err.Error())
    }

    m.Run()

    terminateContainer()
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
