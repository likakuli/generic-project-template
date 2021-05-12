package httputils

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// DefaultHTTPClient will return a default configured http client
var DefaultHTTPClient = func() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			//DisableKeepAlives: true,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}()

// Response is a collection of the response data
type Response struct {
	Status     int
	Header     http.Header
	Body       []byte
	BodyStream io.ReadCloser
}

// Do will excute the request  with the default http client
func Do(req *http.Request) (*Response, error) {
	return ClientDo(DefaultHTTPClient, req)
}

// ClientDo will excute the request with a specific http client
func ClientDo(client *http.Client, req *http.Request, streamResp ...bool) (*Response, error) {
	isStream := false
	if len(streamResp) > 0 && streamResp[0] == true {
		isStream = true
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// return a stream
	if isStream {
		return &Response{
			Status:     resp.StatusCode,
			Header:     resp.Header,
			BodyStream: resp.Body,
		}, nil
	}

	defer resp.Body.Close()
	// return the data
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"read response body failed:%v", err)
	}
	return &Response{
			Status: resp.StatusCode,
			Header: resp.Header,
			Body:   body,
		},
		nil
}

// NewRequest will create a http request with the specified data
func NewRequest(ctx context.Context, method string, url string, headers map[string]string, query url.Values, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	if query != nil {
		q := req.URL.Query()
		for k, vs := range query {
			for _, v := range vs {
				q.Add(k, v)
			}
		}
		req.URL.RawQuery = q.Encode()
	}
	return req, nil
}

func doRequest(ctx context.Context, method string, url string, headers map[string]string, query url.Values, body io.Reader) (*Response, error) {
	req, err := NewRequest(ctx, method, url, headers, query, body)
	if err != nil {
		return nil, err
	}

	return Do(req)
}

// Get will get remote data with custom headers
func Get(ctx context.Context, url string, headers map[string]string, query url.Values) (*Response, error) {
	return doRequest(ctx, "GET", url, headers, query, nil)
}

// Post will create remote resource
func Post(ctx context.Context, url string, headers map[string]string, query url.Values, body io.Reader) (*Response, error) {
	return doRequest(ctx, "POST", url, headers, query, body)
}

// PostForm will create remote resource with x-www-form-urlencoded format data
func PostForm(ctx context.Context, url string, form url.Values) (*Response, error) {
	return Post(ctx,
		url,
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		nil,
		strings.NewReader(form.Encode()))
}

// Put will update a remote resource
func Put(ctx context.Context, url string, headers map[string]string, query url.Values, body io.Reader) (*Response, error) {
	return doRequest(ctx, "PUT", url, headers, query, body)
}

// Patch will partially update a remote resource
func Patch(ctx context.Context, url string, headers map[string]string, query url.Values, body io.Reader) (*Response, error) {
	return doRequest(ctx, "PATCH", url, headers, query, body)
}

// Delete will delete remote resource
func Delete(ctx context.Context, url string, headers map[string]string, query url.Values) (*Response, error) {
	return doRequest(ctx, "DELETE", url, headers, query, nil)
}
