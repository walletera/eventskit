//go:build rabbitmq_client_test

package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/walletera/message-processor/events"
	"github.com/walletera/message-processor/fake"
	"github.com/walletera/message-processor/messages"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestRabbitMQClientWithDefaults(t *testing.T) {
	ctx := context.Background()

	stopRabbitMQ := runRabbitMQ(ctx)
	defer stopRabbitMQ(ctx)

	consumer, err := NewClient()
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer consumer.Close()

	messagesCh, err := consumer.Consume()
	if err != nil {
		t.Errorf("error executing consumer.Consumer: %s", err.Error())
		return
	}

	randomNum := rand.New(rand.NewSource(time.Now().Unix())).Int63()
	data := fmt.Sprintf("a data with a random number %d", randomNum)
	event := fake.Event{
		FakeID:              "FakeID",
		FakeType:            "FakeType",
		FakeCorrelationID:   "FakeCorrelationID",
		FakeDataContentType: "FakeDataContentType",
		FakeData:            data,
	}

	publisher, err := NewClient()
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer publisher.Close()

	err = publisher.Publish(ctx, event, consumer.QueueName())
	if err != nil {
		t.Errorf("error publishing data: %s", err.Error())
		return
	}

	waitForMessageWithTimeout(t, event, messagesCh, 5*time.Second)
}

func TestRabbitMQClientWithOptions(t *testing.T) {
	ctx := context.Background()
	stopRabbitMQ := runRabbitMQ(ctx)

	defer stopRabbitMQ(ctx)

	const testRoutingKey = "test.routing.key"

	consumer, err := NewClient(
		WithQueueName("test-queue"),
		WithExchangeName("test-exchange"),
		WithExchangeType(ExchangeTypeTopic),
		WithConsumerRoutingKeys(testRoutingKey),
	)
	if err != nil {
		t.Errorf("error creating Client: %s", err.Error())
		return
	}
	defer consumer.Close()

	messagesCh, err := consumer.Consume()
	if err != nil {
		t.Errorf("error executing consumer.Consumer: %s", err.Error())
		return
	}

	randomNum := rand.New(rand.NewSource(time.Now().Unix())).Int63()
	data := fmt.Sprintf("a data with a random number %d", randomNum)
	event := fake.Event{
		FakeID:              "FakeID",
		FakeType:            "FakeType",
		FakeCorrelationID:   "FakeCorrelationID",
		FakeDataContentType: "FakeDataContentType",
		FakeData:            data,
	}

	publisher, err := NewClient(
		WithQueueName("test-queue"),
		WithExchangeName("test-exchange"),
		WithExchangeType(ExchangeTypeTopic),
		WithConsumerRoutingKeys(testRoutingKey),
	)
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer publisher.Close()

	err = publisher.Publish(ctx, event, testRoutingKey)
	if err != nil {
		t.Errorf("error publishing message: %s", err.Error())
		return
	}

	waitForMessageWithTimeout(t, event, messagesCh, 5*time.Second)
}

func TestRabbitMQClientWithOptionsButEmptyExchangeType(t *testing.T) {
	ctx := context.Background()
	stopRabbitMQ := runRabbitMQ(ctx)

	defer stopRabbitMQ(ctx)

	const testRoutingKey = "test.routing.key"

	_, err := NewClient(
		WithQueueName("test-queue"),
		WithExchangeName("test-exchange"),
		WithConsumerRoutingKeys(testRoutingKey),
	)
	assert.Error(t, err)
}

func TestRabbitMQClientWithLoad(t *testing.T) {
	ctx := context.Background()
	stopRabbitMQ := runRabbitMQ(ctx)

	defer stopRabbitMQ(ctx)

	const testRoutingKey = "test.routing.key"

	consumer, err := NewClient(
		WithQueueName("test-queue"),
		WithExchangeName("test-exchange"),
		WithExchangeType(ExchangeTypeTopic),
		WithConsumerRoutingKeys(testRoutingKey),
	)
	if err != nil {
		t.Errorf("error creating Client: %s", err.Error())
		return
	}
	defer consumer.Close()

	messagesCh, err := consumer.Consume()
	if err != nil {
		t.Errorf("error executing consumer.Consumer: %s", err.Error())
		return
	}

	randomNum := rand.New(rand.NewSource(time.Now().Unix())).Int63()
	data := fmt.Sprintf("a data with a random number %d", randomNum)
	event := fake.Event{
		FakeID:              "FakeID",
		FakeType:            "FakeType",
		FakeCorrelationID:   "FakeCorrelationID",
		FakeDataContentType: "FakeDataContentType",
		FakeData:            data,
	}

	publisher, err := NewClient(
		WithQueueName("test-queue"),
		WithExchangeName("test-exchange"),
		WithExchangeType("topic"),
		WithConsumerRoutingKeys(testRoutingKey),
	)
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer publisher.Close()

	n := 1000
	for i := 0; i < n; i++ {
		err = publisher.Publish(ctx, event, testRoutingKey)
		if err != nil {
			t.Errorf("error publishing message #%d: %s", i, err.Error())
			return
		}
		time.Sleep(1 * time.Millisecond)
	}

	waitForNMessagesWithTimeout(t, n, messagesCh, 5*time.Second)
}

func runRabbitMQ(ctx context.Context) func(ctx context.Context) {
	req := testcontainers.ContainerRequest{
		Image: "rabbitmq:3.13.3-management",
		Name:  "rabbitmq",
		User:  "rabbitmq",
		ExposedPorts: []string{
			fmt.Sprintf("%d:%d", DefaultPort, DefaultPort),
			fmt.Sprintf("%d:%d", ManagementUIPort, ManagementUIPort),
		},
		WaitingFor: wait.NewExecStrategy([]string{"rabbitmqadmin", "list", "queues"}).WithStartupTimeout(20 * time.Second),
	}
	mockserverC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		fmt.Printf("failed starting rabbitmq container: %s\n", err.Error())
		os.Exit(1)
	}

	return func(ctx context.Context) {
		if err := mockserverC.Terminate(ctx); err != nil {
			fmt.Printf("failed to terminate container: %s\n", err.Error())
		}
	}
}

func waitForMessageWithTimeout(t *testing.T, event events.EventData, messagesCh <-chan messages.Message, duration time.Duration) {
	timeout := time.After(duration)
	select {
	case receivedMessage, ok := <-messagesCh:
		if !ok {
			t.Errorf("channel closed")
			return
		}
		var unmarshalledEvent fake.Event
		err := json.Unmarshal(receivedMessage.Payload(), &unmarshalledEvent)
		if err != nil {
			t.Errorf("failed unmarshalling received message: %s", err.Error())
			return
		}
		assert.Equal(t, event, unmarshalledEvent)
	case <-timeout:
		t.Errorf("timeout waiting message")
	}
}

func waitForNMessagesWithTimeout(t *testing.T, n int, messagesCh <-chan messages.Message, duration time.Duration) {
	timeout := time.After(duration)
	messagesCount := 0
	for {
		select {
		case _, ok := <-messagesCh:
			if !ok {
				t.Errorf("channel closed")
				return
			}
			messagesCount++
			if messagesCount == n {
				return
			}
		case <-timeout:
			t.Errorf("timeout waiting messages")
			return
		}
	}
}
