package main

import (
	"github.com/dbugapp/dbug-go/dbug"
)

// User struct example (mirrors README)
type User struct {
	ID           int
	Name         string
	IsActive     bool
	privateNotes string // Private fields are shown with type info
}

func main() {
	// Example sending a struct
	currentUser := User{
		ID:           123,
		Name:         "Alice",
		IsActive:     true,
		privateNotes: "Needs follow-up",
	}
	// Send the user struct to the Dbug app
	dbug.Go(currentUser)

	// Example sending a map
	myMap := map[string]any{"key": "value", "count": 42}
	// Send the map to the Dbug app
	dbug.Go(myMap)

	// Example sending multiple items (also shown in README)
	message := "Finished examples."
	status := 200
	dbug.Go(message, status)
}
