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

Simply pass any variable to `dbug.Go()`.

See the [examples/main.go](examples/main.go) file for a runnable example.

### Sending Multiple Variables

The `Go` function accepts multiple arguments. Each argument is sent as a separate item to the Dbug app.

See the [examples/main.go](examples/main.go) file for a runnable example.

---

## License

This project is open-sourced under the [MIT license](https://opensource.org/licenses/MIT).

