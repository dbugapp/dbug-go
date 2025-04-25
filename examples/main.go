package main

import (
	"github.com/dbugapp/dbug-go/dbug"
)

func main() {
	basicExample()
	stringExample()
}

func basicExample() {
	dbug.Send(map[string]interface{}{
		"message": "Hello from Go!",
		"user": map[string]interface{}{
			"id":    101,
			"email": "hello@example.com",
		},
	})
}

func stringExample() {
	dbug.Send("Hello from Go!")
}
