package main

import (
	"github.com/dbugapp/dbug-go/dbug"
)

func main() {
	basicExample()
	stringExample()
}

func basicExample() {
	dbug.Go(map[string]interface{}{
		"message": "Hello from Go!",
		"user": map[string]interface{}{
			"id":    101,
			"email": "hello@example.com",
		},
	})
}

func stringExample() {
	dbug.Go("Hello from Go!")
}
