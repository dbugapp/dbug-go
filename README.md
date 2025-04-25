# Dbug Go SDK

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Send debug payloads from Go to the [dbug desktop app](https://github.com/dbugapp/desktop). The `dbug-go` SDK lets you send structured debug information from your Go application to the local `dbug` desktop app for live inspection.

---

## Features

- Serializes Go data structures to JSON
- Handles common serialization errors
- Sends payloads over HTTP to the dbug desktop app
- Customizable endpoint for development flexibility

---

## Installation

```bash
go get github.com/dbugapp/dbug-go
```

---

## Usage

### Basic Example

```go
package main

import (
    "github.com/dbugapp/dbug-go/dbug"
)

func main() {
    dbug.Send(map[string]interface{}{
        "event": "user.registered",
        "user": map[string]interface{}{
            "id":    123,
            "email": "user@example.com",
        },
    })
}
```

This sends the payload to the default dbug server at `http://127.0.0.1:53821`.

---

### Custom Endpoint

You can change the target endpoint using `SetEndpoint()`:

```go
dbug.SetEndpoint("http://localhost:54000")
dbug.Send(map[string]interface{}{
    "event": "order.completed",
    "order": map[string]interface{}{
        "id":     98765,
        "amount": 49.99,
    },
})
```

---


## License

This project is open-sourced under the [MIT license](https://opensource.org/licenses/MIT).

