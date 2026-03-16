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

type Order struct {
	OrderID  string
	Total    float64
	Shipped  bool
}

func main() {
	// Single payload — one cell
	dbug.Go("single string")

	// Multiple payloads — all stacked in one cell
	currentUser := User{
		ID:           123,
		Name:         "Alice",
		IsActive:     true,
		privateNotes: "Needs follow-up",
	}
	order := Order{OrderID: "ORD-999", Total: 49.95, Shipped: false}
	tags := []string{"vip", "trial"}
	meta := map[string]any{"source": "web", "version": 2}

	dbug.Go(currentUser, order, tags, meta)
}
