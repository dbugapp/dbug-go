# Dbug Go Agent

[![Test](https://github.com/dbugapp/dbug-go/actions/workflows/test.yml/badge.svg)](https://github.com/dbugapp/dbug-go/actions/workflows/test.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Send data from your running Go application directly to the [Dbug desktop app](https://github.com/dbugapp/desktop) for live inspection. This package makes it trivial to visualize the state of your Go variables and data structures without interrupting your program.

---

## Features

- **Simple Interface:** Just call `dbug.Go()` with almost any Go variable, struct, map, slice, or other data type to see it in the Dbug app.
- **Variadic Sending:** Send multiple different variables in a single `dbug.Go()` call; each will appear separately.

---

## Installation

```bash
go get github.com/dbugapp/dbug-go
```

---

## Usage

### Sending a Single Variable

Simply pass any variable to `dbug.Go()`:

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
	// Send the user struct to the Dbug app
	dbug.Go(currentUser)

	myMap := map[string]any{"key": "value", "count": 42}
	// Send the map to the Dbug app
	dbug.Go(myMap)
}
```

### Sending Multiple Variables

The `Go` function accepts multiple arguments. Each argument is sent as a separate item to the Dbug app.

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

	// Send user, order, message, and function details
	dbug.Go(user, order, message, calculateTotal)
}
```

---

## License

This project is open-sourced under the [MIT license](https://opensource.org/licenses/MIT).

