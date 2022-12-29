package middleware

import (
	"log"
	"net"
	"time"

	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Logger(address string) gin.HandlerFunc {

	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal("Cannot access centralized logger.")
	}

	logger := logrus.New()
	hook := logrustash.New(conn, logrustash.DefaultFormatter(logrus.Fields{"type": "search_service"}))
	logger.Hooks.Add(hook)

	return func(c *gin.Context) {

		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)

		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		logger.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUri,
		}).Info()
	}
}
