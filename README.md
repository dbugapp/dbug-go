# Dbug Go SDK

[![Go Test](https://github.com/dbugapp/dbug-go/actions/workflows/test.yml/badge.svg)](https://github.com/dbugapp/dbug-go/actions/workflows/test.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Send debug payloads from Go to the [dbug desktop app](https://github.com/dbugapp/desktop). The `dbug-go` SDK lets you send structured debug information from your Go application to the local `dbug` desktop app for live inspection.

---

## Features

- Serializes Go data structures to JSON, aiming to capture maximum detail.
- Handles common serialization errors and circular references.
- Sends multiple payloads in a single call using variadic `Send` function.
- Provides detailed introspection for functions (signature, params, returns) and channels (direction, buffer, length).
- Includes information about private struct fields (name and type).
- Sends payloads over HTTP to the dbug desktop app with a small delay between multiple payloads.
- Customizable endpoint for development flexibility.

---

## Installation

```bash
go get github.com/dbugapp/dbug-go
```

---

## Usage

### Basic Example

Send a single payload:

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

### Sending Multiple Payloads

The `Send` function accepts multiple arguments. Each argument is serialized and sent as a separate payload to the Dbug app.

```go
package main

import (
	"fmt"
	"github.com/dbugapp/dbug-go/dbug"
)

type Order struct {
	ID string
	Items []string
	privateNote string // Will be shown with type info
}

func processOrder(o Order) error {
	fmt.Printf("Processing %s\n", o.ID)
	// ... processing logic ...
	return nil
}

func main() {
	user := map[string]any{"id": 42, "role": "admin"}
	order := Order{ID: "XYZ123", Items: []string{"item1", "item2"}, privateNote: "urgent"}
	message := "Processing complete."

	// Send user, order, message, and function details as separate payloads
	dbug.Send(user, order, message, processOrder)
}
```

This sends four separate payloads to the default dbug server at `http://127.0.0.1:53821`.

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

