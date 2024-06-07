package rtd_test

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/aztechian/nextride-shortcut/internal/rtd"
	"github.com/stretchr/testify/assert"
)

// FileResponder .
type FileResponder func(req *http.Request) *http.Response

// RoundTrip .
func (f FileResponder) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn FileResponder) *http.Client {
	return &http.Client{
		Transport: FileResponder(fn),
	}
}

func readTestFile(filename string) []byte {
	data, _ := os.ReadFile(filename)
	return data
}

func testClient(filename string) *http.Client {
	return NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(readTestFile(filename))),
			Header:     make(http.Header),
		}
	})
}

func TestGet(t *testing.T) {
	// create a test client with a file responder
	client := testClient("../../test/valid-vehicle-rline.json")

	rtd := rtd.Rtd{
		Client: client,
		Line:   "R",
		Stop:   34575,
	}

	// Call the Get function with the desired parameters
	trip, err := rtd.Get()
	if err != nil {
		t.Fatalf("failed to get upcoming trip: %v", err)
	}

	assert.Equal(t, "114798927", trip.TripID)

	// Additional test assertions
	// ...
}
