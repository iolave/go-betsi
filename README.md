# go-betsi

`go-betsi` is a minimalistic and modular web framework for Go, designed to build web applications and APIs with ease. It provides a simple and elegant way to handle HTTP requests, manage configuration, and integrate with other services.

## Features

- **Minimal and Modular:** `go-betsi` is designed to be lightweight and modular, allowing you to use only the components you need.
- **Type-Safe Routing:** Utilizes Go generics to provide type-safe handlers for requests and responses, catching potential bugs at compile time.
- **Built on `chi`:** Leverages the power and flexibility of the popular `chi` router for robust and efficient routing.
- **Simplified Request Handling:** Provides a convenient `AppRequest` struct that simplifies parsing request data and sending responses.
- **Struct Tag-Based Marshalling:** Uses struct tags to automatically marshal and unmarshal request and response bodies, as well as path parameters.
- **Configuration Management:** A simple and flexible configuration system that allows you to configure your application using a `Config` struct.
- **Middleware Support:** Supports `chi` middlewares, allowing you to easily add cross-cutting concerns like logging, authentication, and more.
- **Structured Logging:** Integrates with the [iolave/go-logger](https://github.com/iolave/go-logger) library for structured logging.
- **Tracing:** Integrates with the [iolave/go-tracing](https://github.com/iolave/go-tracing) library for distributed tracing.
- **Centralized Error Handling:** A consistent error handling mechanism using the [iolave/go-errors](https://github.com/iolave/go-errors) library.

## Installation

To install `go-betsi`, use `go get`:

```bash
go get github.com/iolave/go-betsi
```

## Getting Started

Here's a simple example of how to create a "Hello, World!" web server with `go-betsi`:

```go
package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/iolave/go-betsi"
	"github.com/iolave/go-logger"
)

func main() {
	// Create a new logger
	l, err := logger.New(logger.LEVEL_DEBUG, "my-app", "v1.0.0")

	// Create a new router
	r := betsi.NewRouter()

	// Create a new app
	app, err := betsi.New(betsi.Config{
		Logger: l,
		Server: &betsi.ServerConfig{
			Port:   3000,
			Router: r,
		},
	})
	if err != nil {
		panic(err)
	}

    type HandlerRequest struct {
    	ID   string `ar:"path=id" json:"id"`
    	Body struct {
    		Content string `json:"content"`
    	} `ar:"body=json" json:"body"`
    }

	// Define a handler for the "/posts/{id}" route that returns an empty response
	r.Patch("/posts/{id}", goapp.NewHandler(func(ar goapp.AppRequest[HandlerRequest, any]) {
		ctx := ar.Context()
		in, err := ar.ParseRequest()
		if err != nil {
			ar.SendError(ctx, err)
			return
		}
		log.InfoWithData(ctx, "patch_posts", map[string]any{
			"input": in,
		})
		ar.SendJSON(ctx, map[string]any{})
	}))

	// Start the server
	app.Start()
}
```

## Configuration

The `go-betsi` application is configured using the `betsi.Config` struct:

```go
type Config struct {
	// Logger is the logger used to log messages.
	Logger logger.Logger

	// Server is the configuration for the http server.
	// If nil, the http server will not be started.
	Server *ServerConfig
}

type ServerConfig struct {
	// Port is the port for the http server.
	Port int
	// Router is the router used to handle the requests.
	Router *Router
}
```

## Routing

Routing is handled by the `betsi.Router`, which is a wrapper around `chi.Router`. You can define routes for all standard HTTP methods:

```go
r := betsi.NewRouter()

r.Get("/users", getAllUsersHandler)
r.Post("/users", createUserHandler)
r.Get("/users/{id}", getUserHandler)
r.Put("/users/{id}", updateUserHandler)
r.Delete("/users/{id}", deleteUserHandler)
```

### Path Parameters

Path parameters can be defined in the route pattern (e.g., `{id}`) and accessed from the `AppRequest` using struct tags:

```go
type GetUserRequest struct {
    ID string `ar:"path=id"`
}

func getUserHandler(ar betsi.AppRequest[GetUserRequest, any]) {
    // The GetUserRequest struct will be automatically parsed from the request
    req, err := ar.ParseRequest()
    if err != nil {
        ar.SendError(ar.Context(), err)
        return
    }

    // ...
}
```

## Handlers

Handlers are functions that take an `betsi.AppRequest` as an argument. The `AppRequest` provides methods for parsing the request and sending a response.

```go
func (ar AppRequest[In, _]) ParseRequest() (*In, error)
func (ar AppRequest[_, Out]) SendJSON(ctx context.Context, v Out)
func (ar AppRequest[_, _]) SendError(ctx context.Context, err error)
```

## Type-Safe Requests and Responses

`go-betsi` uses generics to provide type-safe handlers. The `In` type parameter is used for the request body, and the `Out` type parameter is for the response body.

```go
type CreateUserRequest struct {
    Body struct {
        Name  string `json:"name"`
        Email string `json:"email"`
    } `ar:"body=json"`
}

type CreateUserResponse struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func createUserHandler(ar betsi.AppRequest[CreateUserRequest, CreateUserResponse]) {
    req, err := ar.ParseRequest()
    if err != nil {
        ar.SendError(ar.Context(), err)
        return
    }

    // ... create user ...

    ar.SendJSON(ar.Context(), CreateUserResponse{
        ID:    "123",
        Name:  req.Body.Name,
        Email: req.Body.Email,
    })
}

// Using it in the Router
func main() {
	r := betsi.NewRouter()
	r.Get("/users", betsi.NewHandler(getUsersHandler))
}
```

## Error Handling

`go-betsi` uses the `github.com/iolave/go-errors` library for error handling. The `AppRequest.SendError` method sends a JSON error response to the client.

```go
ar.SendError(ar.Context(), errors.NewBadRequest("Invalid request body", nil))
```

## Middlewares

Since `go-betsi`'s router is built on `chi`, you can use any `chi`-compatible middleware.

```go
r.Use(middleware.Logger)
r.Use(middleware.Recoverer)
```

You can also use the built-in middlewares from the `pkg/middleware` package.
