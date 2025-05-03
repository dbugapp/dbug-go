package dbug_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dbugapp/dbug-go/dbug"
)

type User struct {
	ID    int
	Email string
}

func TestVariadicSendPayloads(t *testing.T) {
	var receivedBodies [][]byte
	var mu sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		require.NoError(t, err)

		mu.Lock()
		receivedBodies = append(receivedBodies, body)
		mu.Unlock()
	}))
	defer server.Close()

	dbug.SetEndpoint(server.URL)

	payload1 := map[string]interface{}{ "msg": "first" }
	payload2 := User{ID: 456, Email: "another@example.com"}
	payload3 := 12345

	dbug.Go(payload1, payload2, payload3)

	mu.Lock()
	require.Len(t, receivedBodies, 3, "Should have received 3 separate requests")

	var parsed1 map[string]interface{}
	err := json.Unmarshal(receivedBodies[0], &parsed1)
	require.NoError(t, err)
	assert.Equal(t, "first", parsed1["msg"])

	var parsed2 map[string]interface{}
	err = json.Unmarshal(receivedBodies[1], &parsed2)
	require.NoError(t, err)
	// Check struct fields - note the key prefixing
	assert.Equal(t, float64(456), parsed2["dbug_test.User.ID"], "ID should match")
	assert.Equal(t, "another@example.com", parsed2["dbug_test.User.Email"], "Email should match")

	var parsed3 interface{}
	err = json.Unmarshal(receivedBodies[2], &parsed3)
	require.NoError(t, err)
	assert.Equal(t, float64(12345), parsed3, "Third payload should be the integer")
	mu.Unlock()
}

func TestSerializeStructWithCircularReference(t *testing.T) {
	type Node struct {
		Name string
		Next *Node
	}

	a := &Node{Name: "A"}
	b := &Node{Name: "B", Next: a}
	a.Next = b

	jsonBytes, err := dbug.GoTestable(a)
	require.NoError(t, err)

	var parsed map[string]interface{}
	err = json.Unmarshal(jsonBytes, &parsed)
	require.NoError(t, err)

	require.Contains(t, parsed, "dbug_test.Node.Name")
	assert.Equal(t, "A", parsed["dbug_test.Node.Name"])

	require.Contains(t, parsed, "dbug_test.Node.Next")
	nextMap, ok := parsed["dbug_test.Node.Next"].(map[string]interface{})
	require.True(t, ok)

	require.Contains(t, nextMap, "dbug_test.Node.Next")
	assert.Equal(t, "[circular]", nextMap["dbug_test.Node.Next"])
}

func TestSerializeNilAndBasicTypes(t *testing.T) {
	data := map[string]interface{}{
		"nilVal": nil,
		"string": "hello",
		"int":    42,
		"float":  3.14,
		"bool":   true,
		"slice":  []string{"a", "b"},
		"map":    map[string]int{"x": 1},
	}

	jsonBytes, err := dbug.GoTestable(data)
	require.NoError(t, err)

	var parsed map[string]interface{}
	err = json.Unmarshal(jsonBytes, &parsed)
	require.NoError(t, err)

	assert.Nil(t, parsed["nilVal"])
	assert.Equal(t, "hello", parsed["string"])
	assert.Equal(t, float64(42), parsed["int"])
	assert.Equal(t, 3.14, parsed["float"])
	assert.Equal(t, true, parsed["bool"])
	assert.IsType(t, []interface{}{}, parsed["slice"])
	assert.Len(t, parsed["slice"].([]interface{}), 2)
	assert.IsType(t, map[string]interface{}{}, parsed["map"])
	assert.Equal(t, float64(1), parsed["map"].(map[string]interface{})["x"])
}

func TestSerializeFunctionDetailed(t *testing.T) {
	fn := func(i int, s string) (bool, error) { return true, nil }
	jsonBytes, err := dbug.GoTestable(fn)
	require.NoError(t, err)

	var parsed map[string]interface{}
	err = json.Unmarshal(jsonBytes, &parsed)
	require.NoError(t, err)

	signature := "func(int, string) (bool, error)"
	require.Contains(t, parsed, signature)

	details, ok := parsed[signature].(map[string]interface{})
	require.True(t, ok)

	assert.Equal(t, []interface{}{"int", "string"}, details["input_types"])
	assert.Equal(t, []interface{}{"bool", "error"}, details["output_types"])
	assert.False(t, details["is_variadic"].(bool))
}

func TestSerializeChannelDetailed(t *testing.T) {
	chSend := make(chan<- int, 5)
	chRecv := make(<-chan string)
	chBoth := make(chan bool)
	var chNil chan float64

	payload := map[string]interface{}{
		"sendOnly": chSend,
		"recvOnly": chRecv,
		"bothDir":  chBoth,
		"nilChan":  chNil,
	}

	jsonBytes, err := dbug.GoTestable(payload)
	require.NoError(t, err)

	var parsed map[string]interface{}
	err = json.Unmarshal(jsonBytes, &parsed)
	require.NoError(t, err)

	// Send Only
	sendMap, ok := parsed["sendOnly"].(map[string]interface{})
	require.True(t, ok)
	require.Contains(t, sendMap, "chan<- int")
	sendDetails, ok := sendMap["chan<- int"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "int", sendDetails["element_type"])
	assert.Equal(t, "send-only", sendDetails["direction"])
	assert.Equal(t, float64(5), sendDetails["capacity"])
	assert.Equal(t, float64(0), sendDetails["length"])
	assert.False(t, sendDetails["is_nil"].(bool))

	// Recv Only
	recvMap, ok := parsed["recvOnly"].(map[string]interface{})
	require.True(t, ok)
	require.Contains(t, recvMap, "<-chan string")
	recvDetails, ok := recvMap["<-chan string"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "string", recvDetails["element_type"])
	assert.Equal(t, "receive-only", recvDetails["direction"])
	assert.Equal(t, float64(0), recvDetails["capacity"])
	assert.Equal(t, float64(0), recvDetails["length"])
	assert.False(t, recvDetails["is_nil"].(bool))

	// Both Directions
	bothMap, ok := parsed["bothDir"].(map[string]interface{})
	require.True(t, ok)
	require.Contains(t, bothMap, "chan bool")
	bothDetails, ok := bothMap["chan bool"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "bool", bothDetails["element_type"])
	assert.Equal(t, "send-receive", bothDetails["direction"])
	assert.Equal(t, float64(0), bothDetails["capacity"])
	assert.Equal(t, float64(0), bothDetails["length"])
	assert.False(t, bothDetails["is_nil"].(bool))

	// Nil Chan
	nilMap, ok := parsed["nilChan"].(map[string]interface{})
	require.True(t, ok)
	require.Contains(t, nilMap, "chan float64")
	nilDetails, ok := nilMap["chan float64"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "float64", nilDetails["element_type"])
	assert.Equal(t, "send-receive", nilDetails["direction"])
	assert.Equal(t, float64(0), nilDetails["capacity"])
	assert.Equal(t, float64(0), nilDetails["length"])
	assert.True(t, nilDetails["is_nil"].(bool))
}

type StructWithPrivate struct {
	PublicField  string
	privateField int
	privateArr   [5]string
}

func TestSerializeStructPrivateFields(t *testing.T) {
	s := StructWithPrivate{
		PublicField:  "visible",
		privateField: 99,
		privateArr:   [5]string{"a", "b"},
	}

	jsonBytes, err := dbug.GoTestable(s)
	require.NoError(t, err)

	var parsed map[string]interface{}
	err = json.Unmarshal(jsonBytes, &parsed)
	require.NoError(t, err)

	prefix := reflect.TypeOf(s).String() // "dbug_test.StructWithPrivate"

	publicFieldName := fmt.Sprintf("%s.PublicField", prefix)
	privateIntName := fmt.Sprintf("%s.privateField", prefix)
	privateArrName := fmt.Sprintf("%s.privateArr", prefix)

	require.Contains(t, parsed, publicFieldName)
	assert.Equal(t, "visible", parsed[publicFieldName])

	require.Contains(t, parsed, privateIntName)
	assert.Equal(t, "privateField [int]", parsed[privateIntName])

	require.Contains(t, parsed, privateArrName)
	assert.Equal(t, "privateArr [[5]string]", parsed[privateArrName]) // Checking the specific combined name+type format

}
