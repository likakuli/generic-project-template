package httputils

import (
	"net/url"
	"strings"
)

// PackURL will pack the url and the query
func PackURL(addr string, values url.Values) (string, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for k, vs := range values {
		for _, v := range vs {
			q.Add(k, v)
		}
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}

// PackURLPath will pack the url with a url template and args
func PackURLPath(tpl string, args map[string]string) string {
	if args == nil {
		return tpl
	}
	for k, v := range args {
		tpl = strings.Replace(tpl, "{"+k+"}", url.QueryEscape(v), 1)
	}
	return tpl
}
