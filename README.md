# Dbug Go Agent

[![Test](https://github.com/dbugapp/dbug-go/actions/workflows/test.yml/badge.svg)](https://github.com/dbugapp/dbug-go/actions/workflows/test.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Send debug data from your Go application to the [Dbug desktop app](https://github.com/dbugapp/desktop) for live inspection. This package acts as a lightweight agent, making it trivial to visualize the state of your Go variables and data structures.

---

## Features

- **Simple Interface:** Just call `dbug.Send()` with almost any Go variable, struct, map, slice, or other data type.
- **Variadic Sending:** Send multiple different variables in a single `dbug.Send()` call; each will appear separately in the Dbug app.
- **Zero Configuration (Default):** Works out-of-the-box by sending data to the default Dbug desktop app endpoint (`http://127.0.0.1:53821`).

---

## Installation

```bash
go get github.com/dbugapp/dbug-go
```

---

## Usage

### Sending a Single Variable

Simply pass any variable to `dbug.Send()`:

```go
package main

import (
    "github.com/dbugapp/dbug-go/dbug"
)

type User struct {
	ID int
	Name string
	IsActive bool
	privateNotes string // Private fields are shown with type info
}

func main() {
	currentUser := User{
		ID: 123,
		Name: "Alice",
		IsActive: true,
		privateNotes: "Needs follow-up",
	}
	dbug.Send(currentUser)

	myMap := map[string]any{"key": "value", "count": 42}
	dbug.Send(myMap)
}
```

### Sending Multiple Variables

The `Send` function accepts multiple arguments. Each argument is sent as a separate payload to the Dbug app.

```go
package main

import (
	"fmt"
	"github.com/dbugapp/dbug-go/dbug"
)

type Order struct {
	ID string
	Items []string
	privateValue int // Private fields are shown with type info
}

func calculateTotal(o Order) float64 {
	// Dummy calculation
	return float64(len(o.Items) * 10)
}

func main() {
	user := map[string]any{"id": 42, "role": "admin"}
	order := Order{ID: "XYZ123", Items: []string{"item1", "item2"}, privateValue: 99}
	message := "Processing order..."

	// Send user, order, message, and function details as separate payloads
	dbug.Send(user, order, message, calculateTotal)
}
```

This sends four separate payloads to the default Dbug server at `http://127.0.0.1:53821`.

---

## License

This project is open-sourced under the [MIT license](https://opensource.org/licenses/MIT).

