package main

import (
	"github.com/dbugapp/dbug-go/dbug"
)

type User struct {
	ID           int
	Name         string
	IsActive     bool
	privateNotes string 
}

func main() {
	currentUser := User{
		ID:           123,
		Name:         "Alice",
		IsActive:     true,
		privateNotes: "Needs follow-up",
	}
	dbug.Go(currentUser)

	myMap := map[string]any{"key": "value", "count": 42}
	dbug.Go(myMap)

	message := "Finished examples."
	status := 200
	dbug.Go(message, status)
}
