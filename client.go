package firego

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
)

// Client communicates with the Firecracker HTTP API over a Unix domain socket.
// Create one with [New]; a zero Client is not valid.
type Client struct {
	http       *http.Client
	socketPath string
}

// New returns a Client that connects to the Firecracker API at socketPath.
// The socket must exist when individual API calls are made, not at construction time.
func New(socketPath string) *Client {
	transport := &http.Transport{
		DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
			return net.Dial("unix", socketPath)
		},
	}
	return &Client{
		http:       &http.Client{Transport: transport},
		socketPath: socketPath,
	}
}

// APIError is returned when Firecracker responds with an HTTP 4xx or 5xx status.
type APIError struct {
	// StatusCode is the HTTP status code returned by Firecracker.
	StatusCode int
	// FaultMessage is the human-readable description of the error.
	FaultMessage string `json:"fault_message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("firecracker API error %d: %s", e.StatusCode, e.FaultMessage)
}

// Ptr returns a pointer to v. It is useful for setting optional pointer fields
// in request structs without declaring intermediate variables:
//
//	&firego.BootSource{BootArgs: firego.Ptr("console=ttyS0")}
func Ptr[T any](v T) *T { return &v }

// do executes an HTTP request against the Firecracker socket.
// A non-nil body is JSON-encoded and sent as the request body.
// HTTP 4xx/5xx responses are decoded into *APIError and returned as an error.
func (c *Client) do(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, "http://localhost"+path, bodyReader)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		var apiErr APIError
		apiErr.StatusCode = resp.StatusCode
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return nil, &APIError{StatusCode: resp.StatusCode, FaultMessage: resp.Status}
		}
		return nil, &apiErr
	}

	return resp, nil
}

// doJSON calls do and, if out is non-nil, JSON-decodes the response body into it.
func (c *Client) doJSON(ctx context.Context, method, path string, body, out interface{}) error {
	resp, err := c.do(ctx, method, path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}
