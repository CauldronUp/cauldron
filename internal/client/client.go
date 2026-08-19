// Package client talks to a running Cauldron server.
//
// The CLI is a thin shell over the control API rather than a second
// implementation of it. That keeps one behaviour to test, and means anything
// the CLI can do is equally available from a test suite, a Makefile or CI
// without shelling out.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// DefaultBase is where `cauldron serve` listens.
const DefaultBase = "http://127.0.0.1:4600"

// Client is a handle on a running server.
type Client struct {
	base string
	http *http.Client
}

// New returns a client for a base URL. An empty base falls back to the
// CAULDRON_URL environment variable, then to the default.
func New(base string) *Client {
	if base == "" {
		base = os.Getenv("CAULDRON_URL")
	}

	if base == "" {
		base = DefaultBase
	}

	return &Client{
		base: strings.TrimRight(base, "/"),
		http: &http.Client{Timeout: 10 * time.Second},
	}
}

// Base returns the server address.
func (c *Client) Base() string { return c.base }

// ErrUnreachable means no server answered.
type ErrUnreachable struct {
	Base string
	Err  error
}

func (e *ErrUnreachable) Error() string {
	return fmt.Sprintf("no Cauldron server at %s. Start one with 'cauldron serve'", e.Base)
}

func (e *ErrUnreachable) Unwrap() error { return e.Err }

// APIError is an error the server reported.
type APIError struct {
	Status  int
	Message string
}

func (e *APIError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("server returned %d", e.Status)
	}

	return e.Message
}

func (c *Client) do(method, path string, body any, out any) error {
	var reader io.Reader

	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return err
		}

		reader = bytes.NewReader(encoded)
	}

	req, err := http.NewRequest(method, c.base+path, reader)
	if err != nil {
		return err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return &ErrUnreachable{Base: c.base, Err: err}
	}
	defer resp.Body.Close()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return &APIError{Status: resp.StatusCode, Message: messageFrom(payload)}
	}

	if out == nil {
		return nil
	}

	return json.Unmarshal(payload, out)
}

// messageFrom pulls the human-readable message out of an error body so the CLI
// can print the server's own words rather than a status code.
func messageFrom(payload []byte) string {
	var wrapper struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(payload, &wrapper); err == nil && wrapper.Error.Message != "" {
		return wrapper.Error.Message
	}

	return strings.TrimSpace(string(payload))
}

// RecipeStatus describes one mounted sandbox.
type RecipeStatus struct {
	Recipe   string   `json:"recipe"`
	Version  string   `json:"version"`
	Fixture  string   `json:"fixture"`
	Requests int      `json:"requests"`
	Faults   int      `json:"armed_faults"`
	Network  []string `json:"network,omitempty"`
	Webhooks int      `json:"webhooks_sent"`
	Errors   []string `json:"injectable_errors"`
}

// Status is the server's overall state.
type Status struct {
	Time    string         `json:"time"`
	Recipes []RecipeStatus `json:"recipes"`
}

// Status fetches the server state.
func (c *Client) Status() (*Status, error) {
	var out Status

	if err := c.do(http.MethodGet, "/_cauldron/status", nil, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

// Exchange is one logged request.
type Exchange struct {
	Seq      int    `json:"Seq"`
	Method   string `json:"Method"`
	Path     string `json:"Path"`
	Status   int    `json:"Status"`
	Fault    string `json:"Fault"`
	Resource string `json:"Resource"`
	Op       string `json:"Op"`
	Network  string `json:"Network"`
}

// Requests fetches the request log for a recipe.
func (c *Client) Requests(recipe string) ([]Exchange, error) {
	var out struct {
		Requests []Exchange `json:"requests"`
	}

	if err := c.do(http.MethodGet, "/_cauldron/"+url.PathEscape(recipe)+"/requests", nil, &out); err != nil {
		return nil, err
	}

	return out.Requests, nil
}

// Seed loads a fixture.
func (c *Client) Seed(recipe, fixture string) error {
	return c.do(http.MethodPost, "/_cauldron/"+url.PathEscape(recipe)+"/seed", map[string]any{"fixture": fixture}, nil)
}

// Reset returns a recipe, or every recipe when the name is empty, to its
// seeded state.
func (c *Client) Reset(recipe string) error {
	if recipe == "" {
		return c.do(http.MethodPost, "/_cauldron/reset", nil, nil)
	}

	return c.do(http.MethodPost, "/_cauldron/"+url.PathEscape(recipe)+"/reset", nil, nil)
}

// Fault describes a failure to arm.
type Fault struct {
	Error string `json:"error"`
	Count int    `json:"count,omitempty"`
	Every int    `json:"every,omitempty"`
	For   string `json:"for,omitempty"`
	Path  string `json:"path,omitempty"`
}

// Arm installs a fault.
func (c *Client) Arm(recipe string, fault Fault) error {
	return c.do(http.MethodPost, "/_cauldron/"+url.PathEscape(recipe)+"/fault", fault, nil)
}

// Network is a degraded link, armed against one recipe.
//
// Durations are strings so they carry their unit over the wire: "800ms" says
// what it means where 800 would not.
type Network struct {
	Latency     string  `json:"latency,omitempty"`
	Jitter      string  `json:"jitter,omitempty"`
	Bandwidth   int     `json:"bandwidth,omitempty"`
	Timeout     string  `json:"timeout,omitempty"`
	Reset       bool    `json:"reset,omitempty"`
	Limit       int     `json:"limit,omitempty"`
	Slice       int     `json:"slice,omitempty"`
	Probability float64 `json:"probability,omitempty"`
	Count       int     `json:"count,omitempty"`
	For         string  `json:"for,omitempty"`
	Path        string  `json:"path,omitempty"`
	Clear       bool    `json:"clear,omitempty"`
}

// Degrade arms network conditions, and returns the server's description of
// what it armed.
func (c *Client) Degrade(recipe string, network Network) (string, error) {
	var out struct {
		Armed   string `json:"armed"`
		Cleared bool   `json:"cleared"`
	}

	if err := c.do(http.MethodPost, "/_cauldron/"+url.PathEscape(recipe)+"/network", network, &out); err != nil {
		return "", err
	}

	return out.Armed, nil
}

// Emit fires a webhook.
func (c *Client) Emit(recipe, event string, data map[string]any) (string, error) {
	var out struct {
		ID       string `json:"id"`
		Endpoint string `json:"endpoint"`
		Status   int    `json:"status"`
	}

	body := map[string]any{"event": event}
	if data != nil {
		body["data"] = data
	}

	if err := c.do(http.MethodPost, "/_cauldron/"+url.PathEscape(recipe)+"/emit", body, &out); err != nil {
		return "", err
	}

	return out.ID, nil
}

// Subscribe registers a webhook endpoint.
func (c *Client) Subscribe(recipe, endpoint string) error {
	return c.do(http.MethodPost, "/_cauldron/"+url.PathEscape(recipe)+"/subscribe", map[string]any{"url": endpoint}, nil)
}

// Advance moves the shared clock forward and returns the new time.
func (c *Client) Advance(duration string) (string, error) {
	var out struct {
		Now string `json:"now"`
	}

	if err := c.do(http.MethodPost, "/_cauldron/clock/advance", map[string]any{"duration": duration}, &out); err != nil {
		return "", err
	}

	return out.Now, nil
}

// Snapshot fetches the state of every running recipe as raw JSON, so the CLI
// can write it straight to a file without needing to understand the format.
func (c *Client) Snapshot() ([]byte, error) {
	var raw json.RawMessage

	if err := c.do(http.MethodGet, "/_cauldron/snapshot", nil, &raw); err != nil {
		return nil, err
	}

	return raw, nil
}

// Restore applies a snapshot.
func (c *Client) Restore(snapshot []byte) ([]string, error) {
	var body any

	if err := json.Unmarshal(snapshot, &body); err != nil {
		return nil, fmt.Errorf("that file is not a Cauldron snapshot")
	}

	var out struct {
		Restored []string `json:"restored"`
	}

	if err := c.do(http.MethodPost, "/_cauldron/restore", body, &out); err != nil {
		return nil, err
	}

	return out.Restored, nil
}
