// Copyright 2018 Adam S Levy
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package jsonrpc2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
)

// Logger allows custom log types to be used with the Client when
// Client.DebugRequest is true.
type Logger interface {
	Println(...interface{})
	Printf(string, ...interface{})
}

// Client embeds http.Client and provides a convenient way to make JSON-RPC 2.0
// requests.
type Client struct {
	http.Client
	DebugRequest bool
	Log          Logger

	BasicAuth bool
	User      string
	Password  string
	Header    http.Header
}

// Request uses c to make a JSON-RPC 2.0 Request to url with the given method
// and params, and then parses the Response using the provided result, which
// should be a pointer so that it may be populated.
//
// If ctx is not nil, it is added to the http.Request.
//
// If an Error Response is received, then an Error type is returned. Other
// potential errors can result from json.Marshal and params, json.Unmarshal and
// result, http.NewRequest and url, or network errors from c.Do.
//
// A pseudorandom uint between 1 and 5000 is used for the Request.ID.
//
// The "Content-Type":"application/json" header is added to the http.Request,
// and then headers in c.Header are added, which may override the
// "Content-Type".
//
// If c.BasicAuth is true then http.Request.SetBasicAuth(c.User, c.Password) is
// be called.
//
// If c.DebugRequest is true then the Request and Response are printed using
// c.Log. If c.Log == nil, then c.Log = log.New(os.Stderr, "", 0).
func (c *Client) Request(ctx context.Context, url, method string,
	params, result interface{}) error {
	// Generate a random ID for this request.
	reqID := rand.Int()%5000 + 1

	// Marshal the JSON RPC Request.
	reqJrpc := Request{ID: reqID, Method: method, Params: params}
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
	reqHTTP, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if ctx != nil {
		reqHTTP = reqHTTP.WithContext(ctx)
	}
	reqHTTP.Header.Add(http.CanonicalHeaderKey("Content-Type"), "application/json")
	for k, v := range c.Header {
		reqHTTP.Header[http.CanonicalHeaderKey(k)] = v
	}
	if c.BasicAuth {
		reqHTTP.SetBasicAuth(c.User, c.Password)
	}

	// Make the request.
	res, err := c.Do(reqHTTP)
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
		return err
	}

	// Unmarshal the HTTP response into a JSON RPC response.
	var resID int
	resJrpc := Response{Result: result, ID: &resID}
	if err := json.Unmarshal(resBytes, &resJrpc); err != nil {
		return err
	}
	if c.DebugRequest {
		if resJrpc.HasError() {
			resJrpc.Result = nil
		}
		fmt.Println("<--", string(resBytes))
		fmt.Println()
	}
	if resJrpc.HasError() {
		return resJrpc.Error
	}
	if resID != reqID {
		return fmt.Errorf("request/response ID mismatch")
	}
	return nil
}
