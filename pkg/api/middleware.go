package api

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func LoggingMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		startTime := time.Now()

		// Call the next handler in the chain
		next(ctx)

		// Log information about the request using logrus
		endTime := time.Now()
		duration := endTime.Sub(startTime)
		logrus.WithFields(logrus.Fields{
			"Time":       endTime.Format("2006-01-02 15:04:05"),
			"RemoteAddr": ctx.RemoteAddr(),
			"Method":     string(ctx.Method()),
			"RequestURI": ctx.URI().String(),
			"Duration":   duration,
		}).Info("Request handled")
	})
}
