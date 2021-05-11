package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/leopoldxx/go-utils/trace"
)

type key int

const (
	tracerLogHandlerID key = 32702 // random key
	realIPValueID      key = 16221
)

// HandleFunc wrap a trace handle func outer the original http handle func
func Trace(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tracer trace.Trace
		if id := c.Request.Header.Get("x-request-id"); len(id) > 0 {
			tracer = trace.WithID(name, id)
		} else {
			tracer = trace.New(name)
		}
		c.Writer.Header().Set("x-request-id", tracer.ID())

		lastRoute, ip := func(r *http.Request) (string, string) {
			lastRoute := strings.Split(r.RemoteAddr, ":")[0]
			if ip, exists := r.Header["X-Real-IP"]; exists && len(ip) > 0 {
				return lastRoute, ip[0]
			}
			if ips, exists := r.Header["X-Forwarded-For"]; exists && len(ips) > 0 {
				return lastRoute, ips[0]
			}
			return lastRoute, lastRoute
		}(c.Request)

		tracer.Infof("event=[request-in] remote=[%s] route=[%s] method=[%s] url=[%s]", ip, lastRoute, c.Request.Method, c.Request.URL.String())
		defer tracer.Info("event=[request-out]")

		ctx := context.WithValue(c.Request.Context(), tracerLogHandlerID, tracer)
		ctx = context.WithValue(ctx, realIPValueID, ip)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
