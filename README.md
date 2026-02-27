# Eventskit
[![CI](https://github.com/walletera/eventskit/actions/workflows/ci.yml/badge.svg)](https://github.com/walletera/eventskit/actions/workflows/ci.yml)

**Eventskit** is a modular library for building robust, event-driven systems in Go. It is designed to be both lightweight and highly extensible, making it a great fit for microservices and distributed architectures.

## Key Features

### Generic and Strongly Typed Message Processor

At the heart of Eventskit is a generic, strongly typed message processor, centered around the `Processor` type in the `messages` package. This processor enforces several best practices:

- **Decoupled Architecture:** Cleanly separates message broker integration, deserialization, and message handling, improving maintainability and extensibility.
- **Abstract Error Handling:** Delegates error handling to an abstract acknowledger, making it easy to implement message replay, dead-letter queue parking, or discard strategies when errors occur during processing.
- **Type Safety via Generics:** Leverages Go generics to enforce compile-time type safety, eliminating the risks and maintenance issues associated with reflection or dynamic typing.

### Built-in Consumers

Eventskit includes ready-to-use consumers for popular platforms:

- **RabbitMQ**: Integrate seamlessly with RabbitMQ for pub/sub and queue-based messaging.
- **EventStoreDB**: Process event streams using EventStoreDB, enabling event-sourcing architectures.
- **Webhooks**: Easily consume events delivered via webhooks.

### Event Sourcing Support

Native utilities and abstractions make it straightforward to apply event sourcing patterns, track system state through event streams, and implement audit or rebuild mechanisms.

## Getting Started

1. **Install**
```textmate
go get github.com/walletera/eventskit
```

2. **Import & Use**  
   Import only the packages you need and start wiring up processors and consumers.

3. **Extensibility**  
   Eventskit’s modular design allows you to add new brokers, acknowledgers, or processors with minimal effort.

## Contributing

We welcome contributions! Please open issues or submit pull requests to help improve the library.

## License

Distributed under the MIT License. See [LICENSE](LICENSE) for details.

---

**Eventskit**—Build event-driven Go services with clean, type-safe, and modular patterns.
