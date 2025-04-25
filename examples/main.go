package main

import (
	"github.com/dbugapp/dbug-go/dbug"
)

func main() {
	dbug.Send(map[string]interface{}{
		"message": "Hello from Go!",
		"user": map[string]interface{}{
			"id":    101,
			"email": "hello@example.com",
		},
	})
}
