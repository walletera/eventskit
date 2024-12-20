package webhook

import (
    "bytes"
    "context"
    "log/slog"
    "net/http"
    "os"
    "strings"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/walletera/eventskit/messages"
    "github.com/walletera/werrors"
)

const (
    payloadSucceed                        = `succeed`
    payloadFailedWithUnprocessableMessage = `failUnprocessableMessage`
    payloadFailedWithInternalError        = `failInternalError`
)

func TestServer(t *testing.T) {
    var tests = []struct {
        name                string
        httpMethod          string
        notificationPayload []byte
        expectedStatusCode  int
    }{
        {
            name:                "process request successfully",
            httpMethod:          "POST",
            notificationPayload: []byte(payloadSucceed),
            expectedStatusCode:  http.StatusCreated,
        },
        {
            name:                "wrong http method produce MethodNotAllowed http status code",
            httpMethod:          "GET",
            notificationPayload: []byte(payloadSucceed),
            expectedStatusCode:  http.StatusMethodNotAllowed,
        },
        {
            name:                "empty body produce BadRequest error",
            httpMethod:          "POST",
            notificationPayload: nil,
            expectedStatusCode:  http.StatusBadRequest,
        },
        {
            name:                "nack with status code UnprocessableMessage produce BadRequest http status code",
            httpMethod:          "POST",
            notificationPayload: []byte(payloadFailedWithUnprocessableMessage),
            expectedStatusCode:  http.StatusBadRequest,
        },
        {
            name:                "nack with status code InternalError produce InternalError http status code",
            httpMethod:          "POST",
            notificationPayload: []byte(payloadFailedWithInternalError),
            expectedStatusCode:  http.StatusInternalServerError,
        },
    }
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            executeTest(t, test.httpMethod, test.notificationPayload, test.expectedStatusCode)
        })
    }
}

func executeTest(t *testing.T, httpMethod string, requestBody []byte, expectedStatusCode int) {
    server := NewServer(8282, WithLogger(slog.New(slog.NewTextHandler(os.Stdout, nil))))
    msgCh, err := server.Consume()
    require.NoError(t, err)

    ctx, cancelCtx := context.WithTimeout(context.Background(), 2*time.Second)
    go func() {
        for msg := range msgCh {
            assert.NotNil(t, msg.Payload)
            var err error
            switch string(msg.Payload()) {
            case payloadSucceed:
                err = msg.Acknowledger().Ack()
            case payloadFailedWithUnprocessableMessage:
                err = msg.Acknowledger().Nack(messages.NackOpts{
                    Requeue:      false,
                    ErrorCode:    werrors.UnprocessableMessageErrorCode,
                    ErrorMessage: "",
                })
            case payloadFailedWithInternalError:
                err = msg.Acknowledger().Nack(messages.NackOpts{
                    Requeue:      false,
                    ErrorCode:    werrors.InternalErrorCode,
                    ErrorMessage: "",
                })
            }

            assert.NoError(t, err)
        }
        cancelCtx()
    }()

    httpReq, err := http.NewRequest(httpMethod, "http://127.0.0.1:8282/webhooks", bytes.NewReader(requestBody))
    require.NoError(t, err)

    resp, err := http.DefaultClient.Do(httpReq)
    require.NoError(t, err)

    assert.Equal(t, expectedStatusCode, resp.StatusCode)

    err = server.Close()
    require.NoError(t, err)

    <-ctx.Done()
    err = ctx.Err()
    assert.ErrorIs(t, err, context.Canceled)
}

func TestShutdown(t *testing.T) {
    server := NewServer(8282, WithLogger(slog.New(slog.NewTextHandler(os.Stdout, nil))))
    msgCh, err := server.Consume()
    require.NoError(t, err)

    messageProcessingDelay := 10 * time.Millisecond

    ctx, cancelCtx := context.WithTimeout(context.Background(), 2*time.Second)
    go func() {
        for msg := range msgCh {
            assert.NotNil(t, msg.Payload)
            time.Sleep(messageProcessingDelay)
            err := msg.Acknowledger().Ack()
            assert.NoError(t, err)
        }
        cancelCtx()
    }()

    time.AfterFunc(2*time.Millisecond, func() {
        err := server.Close()
        assert.NoError(t, err)
    })

    start := time.Now()
    httpReq, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8282/webhooks", strings.NewReader("some payload"))
    require.NoError(t, err)

    resp, err := http.DefaultClient.Do(httpReq)
    require.NoError(t, err)
    elapsed := time.Since(start)

    assert.Equal(t, http.StatusCreated, resp.StatusCode)
    assert.GreaterOrEqual(t, elapsed, messageProcessingDelay)

    err = server.Close()
    require.NoError(t, err)

    <-ctx.Done()
    assert.ErrorIs(t, ctx.Err(), context.Canceled)
}
