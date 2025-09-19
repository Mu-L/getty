# getty [中文](./README_CN.md)

 *a netty like asynchronous network I/O library*

[![Build Status](https://travis-ci.org/AlexStocks/getty.svg?branch=master)](https://travis-ci.org/AlexStocks/getty)
[![codecov](https://codecov.io/gh/AlexStocks/getty/branch/master/graph/badge.svg)](https://codecov.io/gh/AlexStocks/getty)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/AlexStocks/getty?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/AlexStocks/getty)](https://goreportcard.com/report/github.com/AlexStocks/getty)
![license](https://img.shields.io/badge/license-Apache--2.0-green.svg)

## INTRO

Getty is an asynchronous network I/O library developed in Golang. It operates on TCP, UDP, and WebSocket network protocols, providing a consistent interface [EventListener](https://github.com/AlexStocks/getty/blob/01184614ef72d0cf2dd11894ab31e0dace066b6c/transport/getty.go#L68).

Within Getty, each connection (session) involves two separate goroutines. One handles the reading of TCP streams, UDP packets, or WebSocket packages, while the other manages the logic processing and writes responses into the network write buffer. If your logic processing might take a considerable amount of time, it's recommended to start a new logic process goroutine yourself within codec.go's (Codec)OnMessage method.

Additionally, you can manage heartbeat logic within the (Codec)OnCron method in codec.go. If you're using TCP or UDP, you should send heartbeat packages yourself and then call session.go's (Session)UpdateActive method to update the session's activity timestamp. You can verify if a TCP session has timed out or not in codec.go's (Codec)OnCron method using session.go's (Session)GetActive method.

If you're using WebSocket, you don't need to worry about heartbeat request/response, as Getty handles this task within session.go's (Session)handleLoop method by sending and receiving WebSocket ping/pong frames. Your responsibility is to check whether the WebSocket session has timed out or not within codec.go's (Codec)OnCron method using session.go's (Session)GetActive method.

For code examples, you can refer to [getty-examples](https://github.com/AlexStocks/getty-examples).

## Callback System

Getty provides a robust callback system that allows you to register and manage callback functions for session lifecycle events. This is particularly useful for cleanup operations, resource management, and custom event handling.

### Key Features

- **Thread-safe operations**: All callback operations are protected by mutex locks
- **Replace semantics**: Adding with the same (handler, key) replaces the existing callback in place (position preserved)
- **Panic safety**: During session close, callbacks run in a dedicated goroutine with defer/recover; panics are logged with stack traces and do not escape the close path
- **Ordered execution**: Callbacks are executed in the order they were added

### Usage Example

```go
// Add a close callback
session.AddCloseCallback("cleanup", "resources", func() {
    // Cleanup resources when session closes
    cleanupResources()
})

// Remove a specific callback
// Safe to call even if the pair was never added (no-op)
session.RemoveCloseCallback("cleanup", "resources")

// Callbacks are automatically executed when the session closes
```

**Note**: During session shutdown, callbacks are executed sequentially in a dedicated goroutine to preserve add-order, with defer/recover to log panics without letting them escape the close path.

### Callback Management

- **AddCloseCallback**: Register a callback to be executed when the session closes
- **RemoveCloseCallback**: Remove a previously registered callback (no-op if not found; safe to call multiple times)
- **Thread Safety**: All operations are thread-safe and can be called concurrently

### Type Requirements

The `handler` and `key` parameters must be **comparable types** that support the `==` operator:

**✅ Supported types:**
- **Basic types**: `string`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `uintptr`, `float32`, `float64`, `bool`, `complex64`, `complex128`
  - ⚠️ Avoid `float*`/`complex*` as keys due to NaN and precision semantics; prefer strings/ints
- **Pointer types**: Pointers to any type (e.g., `*int`, `*string`, `*MyStruct`)
- **Interface types**: Interface types are comparable only when their dynamic values are comparable types; using "==" with non-comparable dynamic values will be safely ignored with error log
- **Channel types**: Channel types (compared by channel identity)
- **Array types**: Arrays of comparable elements (e.g., `[3]int`, `[2]string`)
- **Struct types**: Structs where all fields are comparable types

**⚠️ Non-comparable types (will be safely ignored with error log):**
- `map` types (e.g., `map[string]int`)
- `slice` types (e.g., `[]int`, `[]string`)
- `func` types (e.g., `func()`, `func(int) string`)
- Structs containing non-comparable fields (maps, slices, functions)

**Examples:**
```go
// ✅ Valid usage
session.AddCloseCallback("user", "cleanup", callback)
session.AddCloseCallback(123, "cleanup", callback)
session.AddCloseCallback(true, false, callback)

// ⚠️ Non-comparable types (safely ignored with error log)
session.AddCloseCallback(map[string]int{"a": 1}, "key", callback)  // Logged and ignored
session.AddCloseCallback([]int{1, 2, 3}, "key", callback)          // Logged and ignored
```

## About network transmission in getty

In network communication, the data transmission interface of getty does not guarantee that data will be sent successfully; it lacks an internal retry mechanism. Instead, getty delegates the outcome of data transmission to the underlying operating system mechanism. Under this mechanism, if data is successfully transmitted, it is considered a success; if transmission fails, it is regarded as a failure. These outcomes are then communicated back to the upper-layer caller.

Upper-layer callers need to determine whether to incorporate a retry mechanism based on these outcomes. This implies that when data transmission fails, upper-layer callers must handle the situation differently depending on the circumstances. For instance, if the failure is due to a disconnect in the connection, upper-layer callers can attempt to resend the data based on the result of getty's automatic reconnection. Alternatively, if the failure is caused by the sending buffer of the underlying operating system being full, the sender can implement its own retry mechanism to wait for the sending buffer to become available before attempting another transmission.

In summary, the data transmission interface of getty does not come with an inherent retry mechanism; instead, it is up to upper-layer callers to decide whether to implement retry logic based on specific situations. This design approach provides developers with greater flexibility in controlling the behavior of data transmission.

## LICENCE

Apache License 2.0

