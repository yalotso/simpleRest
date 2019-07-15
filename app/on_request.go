package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func OnRequest(logger Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		rs := newRequestScope(logger, time.Now())
		c.Set("Context", rs)

		c.Next()

		logAccess(c, rs)
		logError(c, rs)
	}
}

func GetRequestScope(c *gin.Context) RequestScope {
	return c.MustGet("Context").(RequestScope)
}

func logAccess(c *gin.Context, rs RequestScope) {
	elapsed := float64(time.Now().Sub(rs.Now()).Nanoseconds()) / 1e6
	requestLine := fmt.Sprintf("%s %s %s", c.Request.Method, c.Request.URL.Path, c.Request.Proto)
	rs.Infof("[%.3fms] %s %s %d %d", elapsed, c.ClientIP(), requestLine, c.Writer.Status(), c.Writer.Size())
}

func logError(c *gin.Context, rs RequestScope) {
	if len(c.Errors) == 0 {
		return
	}
	rs.Errorf("%s", c.Errors)
}
