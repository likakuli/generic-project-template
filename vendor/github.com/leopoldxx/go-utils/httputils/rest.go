package httputils

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/leopoldxx/go-utils/trace"
)

// DebugLevel of the debug logs
type DebugLevel int

// Prededfined debuglog level
const (
	Debug0 DebugLevel = 0
	Debug1 DebugLevel = 1
	Debug2 DebugLevel = 2
)

// RestCli represents a restful http request
type RestCli struct {
	cli      *http.Client
	ctx      context.Context
	method   string
	api      string
	host     string
	resource string
	headers  map[string]string
	querys   url.Values
	formData url.Values
	body     io.Reader
	object   interface{}
	into     map[string]interface{}
	debug    DebugLevel
	isStream bool
}

var defaultHTTPClient = func() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}()

// NewRestCli will create a new RestCli object
func NewRestCli() *RestCli {
	return &RestCli{
		cli:      defaultHTTPClient,
		method:   "GET",
		api:      "http://127.0.0.1/",
		host:     "http://127.0.0.1",
		resource: "/",
		headers:  map[string]string{},
		querys:   url.Values{},
		into:     map[string]interface{}{},
		debug:    Debug0,
		isStream: false,
	}
}

// Client will set a user specified http client for the rest request
func (rest *RestCli) Client(cli *http.Client) *RestCli {
	rest.cli = cli
	return rest
}

// Context will set a user specified context for the rest request
func (rest *RestCli) Context(ctx context.Context) *RestCli {
	rest.ctx = ctx
	return rest
}

// Method will set the http method for the rest request
func (rest *RestCli) Method(method string) *RestCli {
	rest.method = method
	return rest
}

// Get will get the rest request
func (rest *RestCli) Get() *RestCli {
	rest.method = "GET"
	return rest
}

// Post will post the rest request
func (rest *RestCli) Post() *RestCli {
	rest.method = "POST"
	return rest
}

// Put will put the rest request
func (rest *RestCli) Put() *RestCli {
	rest.method = "PUT"
	return rest
}

// Delete will delete the rest request
func (rest *RestCli) Delete() *RestCli {
	rest.method = "DELETE"
	return rest
}

// Patch will post the rest request
func (rest *RestCli) Patch() *RestCli {
	rest.method = "PATCH"
	return rest
}

// Host will set the remote host address for the rest request
func (rest *RestCli) Host(host string) *RestCli {
	rest.host = host
	rest.api = rest.host + rest.resource
	return rest
}

// ResourcePath will set the remote api Resource for the rest request
func (rest *RestCli) ResourcePath(resource string) *RestCli {
	rest.resource = resource
	rest.api = rest.host + rest.resource
	return rest
}

// SetHeader will set a header pair for the rest request
func (rest *RestCli) SetHeader(header, value string) *RestCli {
	rest.headers[header] = value
	return rest
}

// ClearHeader will clear a header from the rest request
func (rest *RestCli) ClearHeader(header string) *RestCli {
	delete(rest.headers, header)
	return rest
}

// SetQuery will set a query for the rest request
func (rest *RestCli) SetQuery(query string, value ...string) *RestCli {
	values := []string{}
	for _, v := range value {
		values = append(values, v)
	}
	rest.querys[query] = values
	return rest
}

// ClearQuery will clear a query from the rest request
func (rest *RestCli) ClearQuery(query string) *RestCli {
	delete(rest.querys, query)
	return rest
}

// FormData will set the body object for the form request
func (rest *RestCli) FormData(data url.Values) *RestCli {
	rest.formData = data
	return rest
}

// Object will set the body object for the rest request
func (rest *RestCli) Object(body interface{}) *RestCli {
	rest.object = body
	return rest
}

// Body will set body for the rest request
func (rest *RestCli) Body(body io.Reader) *RestCli {
	rest.body = body
	return rest
}

// Into will store the ok response of the rest request
func (rest *RestCli) Into(status string, resp interface{}) *RestCli {
	rest.into[status] = resp
	return rest
}

// Stream will turn on or turn off the debug process
func (rest *RestCli) Stream() *RestCli {
	rest.isStream = true
	return rest
}

// Debug will turn on or turn off the debug process
func (rest *RestCli) Debug(level ...DebugLevel) *RestCli {
	if len(level) > 0 {
		rest.debug = level[0]
	} else {
		rest.debug = Debug1
	}
	return rest
}

// Do will send the rest request to remote api and process the resp
func (rest *RestCli) Do() (*Response, error) {
	if rest.ctx == nil {
		rest.ctx = context.TODO()
	}

	tracer := trace.GetTraceFromContext(rest.ctx)
	if _, exists := rest.headers["x-request-id"]; !exists {
		rest.headers["x-request-id"] = tracer.ID()
	}
	if rest.debug >= Debug1 {
		tracer.Infof("api: %v, query: %v", rest.api, rest.querys)
	}

	bodyReader := rest.body
	if rest.object != nil {
		body, err := json.Marshal(rest.object)
		if err != nil {
			if rest.debug >= Debug1 {
				tracer.Error("marshal request body failed:", err)
			}
			return nil, err
		}
		bodyReader = bytes.NewReader(body)
		if rest.method == "GET" || rest.method == "HEAD" || rest.method == "DELETE" {
			rest.method = "POST"
		}
		if _, ok := rest.headers["Content-Type"]; !ok {
			rest.headers["Content-Type"] = "application/json"
		}
		if rest.debug >= Debug2 {
			tracer.Infof("req header: %v", rest.headers)
			tracer.Infof("req body: %v", string(body))
		}
	} else if rest.formData != nil {
		body := rest.formData.Encode()
		bodyReader = strings.NewReader(body)
		if rest.method == "GET" || rest.method == "HEAD" || rest.method == "DELETE" {
			rest.method = "POST"
		}
		if _, ok := rest.headers["Content-Type"]; !ok {
			rest.headers["Content-Type"] = "application/x-www-form-urlencoded"
		}
		if rest.debug >= Debug2 {
			tracer.Infof("req header: %v", rest.headers)
			tracer.Infof("req body: %v", body)
		}
	}

	req, err := NewRequest(
		rest.ctx,
		rest.method,
		rest.api,
		rest.headers,
		rest.querys,
		bodyReader)
	if err != nil {
		if rest.debug >= Debug1 {
			tracer.Error("create request failed:", err)
		}
		return nil, err
	}

	if strings.EqualFold(rest.querys.Get("Connection"), "close") {
		req.Close = true
	}

	resp, err := ClientDo(rest.cli, req, true) // always return  a Body Reader, avoid memory copy
	if err != nil {
		if rest.debug >= Debug1 {
			tracer.Error("do request failed:", err)
		}
		return nil, err
	}
	if rest.debug >= Debug1 {
		tracer.Infof("resp status: %v, header: %v", resp.Status, resp.Header)
	}

	if rest.isStream {
		return resp, nil
	}
	defer resp.BodyStream.Close()

	if len(rest.into) > 0 {
		status := strconv.Itoa(resp.Status)
		if len(status) == 3 {
			ss := []string{
				status,            // YYY, 200, 201, 400, 500
				status[:2] + "x",  // YYx, 20x, 30x, 40x, 50x
				status[:1] + "xx", // Yxx, 2xx, 3xx, 4xx, 5xx
				"xxx",             // xxx
			}
			for _, status := range ss {
				if rsp, exist := rest.into[status]; exist {
					err = json.NewDecoder(resp.BodyStream).Decode(rsp)
					if err != nil {
						if rest.debug >= Debug1 {
							tracer.Errorf("unmarshal resp failed: %s", err)
						}
						return nil, err
					}
					if rest.debug >= Debug2 {
						tracer.Infof("resp: %+v", rsp)
					}
					return resp, nil
				}
			}
		}
	}
	return resp, nil
}
