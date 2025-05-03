# go-app

A simple web framework for Go.

## Features

- Dependency injection
- Middlewares
- Request ID
- Request logging
- Request validation
- Vault integration

## Environment Variables

- `VAULT_ADDR`: Vault address (format: `scheme://host:port`).
- `VAULT_TOKEN`: Vault token.
- `CF_ACCESS_CLIENT_ID`: Cloudflare access client id.
- `CF_ACCESS_CLIENT_SECRET`: Cloudflare access client secret.

## Getting Started

### Installation

```bash
go get github.com/pingolabscl/go-app
```

### Usage

```go
package main

import (
        "context"
        "fmt"
        "net/http"

        "github.com/pingolabscl/go-app"
        "github.com/pingolabscl/go-app/logger"
)

func main() {
    app, err := goapp.New(goapp.Config{
        Name: "go-app",
        LogLevel: logger.LEVEL_INFO,
        Port: 3000,
        InsecureSkipVerify: true,
        Vault: goapp.VaultConfig{
            Addr: "https://vault.pingolabs.cl:443",
            Token: "token",
        },
    })
    if err != nil {
        panic(err)
    }

    app.Get("/hello", func(ar goapp.AppRequest) {
        app.Logger.InfoWithData(ar.Context(), "hello", map[string]any{
            "name": "world",
        })
        ar.SendJSON(map[string]any{
            "message": "Hello, world!",
        })
    })

    app.Start()
}
```
