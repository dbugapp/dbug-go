package dbug_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dbugapp/dbug-go/dbug"
)

type User struct {
	ID    int
	Email string
}

func TestSendPayload(t *testing.T) {
	var receivedBody []byte

	// Set up a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		assert.NoError(t, err)
		receivedBody = body
	}))
	defer server.Close()

	dbug.SetEndpoint(server.URL)

	// Send a sample payload
	dbug.Send(map[string]interface{}{
		"event": "user.created",
		"user": User{
			ID:    123,
			Email: "test@example.com",
		},
	})

	// Check that we got something
	assert.NotNil(t, receivedBody)

	var parsed map[string]interface{}
	err := json.Unmarshal(receivedBody, &parsed)
	assert.NoError(t, err)
	assert.Equal(t, "user.created", parsed["event"])
}

func TestSerializeStructWithCircularReference(t *testing.T) {
	type Node struct {
		Name string
		Next *Node
	}

	a := &Node{Name: "A"}
	b := &Node{Name: "B", Next: a}
	a.Next = b // circular reference

	// Should not panic or error
	jsonBytes, err := dbug.SendTestable(a)
	assert.NoError(t, err)

	var parsed map[string]interface{}
	err = json.Unmarshal(jsonBytes, &parsed)
	assert.NoError(t, err)

	assert.Equal(t, "A", parsed["Name"])
	assert.Equal(t, "[circular]", parsed["Next"].(map[string]interface{})["Next"])
}

func TestSerializeNilAndBasicTypes(t *testing.T) {
	data := map[string]interface{}{
		"nilVal":   nil,
		"string":   "hello",
		"int":      42,
		"float":    3.14,
		"bool":     true,
		"slice":    []string{"a", "b"},
		"map":      map[string]int{"a": 1},
		"function": func() {},
	}

	jsonBytes, err := dbug.SendTestable(data)
	assert.NoError(t, err)

	var parsed map[string]interface{}
	err = json.Unmarshal(jsonBytes, &parsed)
	assert.NoError(t, err)
	assert.Equal(t, float64(42), parsed["int"])
	assert.Equal(t, true, parsed["bool"])
}
