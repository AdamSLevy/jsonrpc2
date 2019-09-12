package jsonrpc2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
)

type Logger interface {
	Println(...interface{})
}

type RequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client embeds http.Client and provides a convenient way to make JSON-RPC
// requests.
type Client struct {
	RequestDoer
	DebugRequest bool
	Log          Logger

	BasicAuth bool
	User      string
	Password  string
	Header    http.Header
}

// Request uses c to make a JSON-RPC 2.0 Request to url with the given method
// and params, and then parses the Response using the provided result for
// Response.Result. Thus, result must be a pointer in order for json.Unmarshal
// to populate it. If Request returns nil, then the request and RPC method call
// were successful and result will be populated, if applicable. If the request
// is successful but the RPC method returns an Error Response, then Request
// will return the Error, which can be checked for by attempting a type
// assertion on the returned error.
//
// Request uses a pseudorandom uint32 for the Request.ID.
//
// Requests will have the "Content-Type":"application/json" header added.
//
// Any populated c.Header will then be added to the http.Request, so you may
// override the "Content-Type" header with your own.
//
// If c.BasicAuth is true then http.Request.SetBasicAuth(c.User, c.Password)
// will be called. This will override the same header in c.Header.
//
// If c.DebugRequest is true then the Request and Response will be printed to
// stdout.
func (c *Client) Request(url, method string, params, result interface{}) error {
	// Generate a random ID for this request.
	reqID := rand.Uint32()%200 + 500

	// Marshal the JSON RPC Request.
	reqJrpc := NewRequest(method, reqID, params)
	if c.DebugRequest {
		if c.Log == nil {
			c.Log = log.New(os.Stderr, "", 0)
		}
		c.Log.Println(reqJrpc)
	}
	reqBytes, err := reqJrpc.MarshalJSON()
	if err != nil {
		return err
	}

	// Compose the HTTP request.
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	req.Header.Add(http.CanonicalHeaderKey("Content-Type"), "application/json")
	for k, v := range c.Header {
		req.Header[http.CanonicalHeaderKey(k)] = v
	}
	if c.BasicAuth {
		req.SetBasicAuth(c.User, c.Password)
	}

	// Make the request.
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusBadRequest {
		return fmt.Errorf("http: %v", res.Status)
	}

	// Read the HTTP response.
	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll(http.Response.Body): %v", err)
	}

	// Unmarshal the HTTP response into a JSON RPC response.
	var resID uint32
	resJrpc := Response{Result: result, ID: &resID}
	if err := json.Unmarshal(resBytes, &resJrpc); err != nil {
		return fmt.Errorf("json.Unmarshal(%v): %v", string(resBytes), err)
	}
	if c.DebugRequest {
		if resJrpc.Error != nil {
			resJrpc.Result = nil
		}
		fmt.Println("<--", string(resBytes))
		fmt.Println()
	}
	if resJrpc.Error != nil {
		return *resJrpc.Error
	}
	if resID != reqID {
		return fmt.Errorf("request/response ID mismatch")
	}
	return nil
}


func NewClient(doer RequestDoer) *Client {
	if doer == nil {
		doer = &http.Client{}
	}
	return &Client{RequestDoer: doer}
}
