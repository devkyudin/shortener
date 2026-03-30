package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/devkyudin/shortener/internal/logger"
)

type LoggingMiddleware struct {
	logContainer *logger.Container
}

type (
	responseData struct {
		statusCode int
		size       int
	}

	loggingResponseWriter struct {
		responseWriter http.ResponseWriter
		data           *responseData
	}
)

func (r *loggingResponseWriter) Header() http.Header {
	return r.responseWriter.Header()
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.responseWriter.Write(b)
	r.data.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.responseWriter.WriteHeader(statusCode)
	r.data.statusCode = statusCode
}

func NewLoggingMiddleware(logContainer *logger.Container) *LoggingMiddleware {
	return &LoggingMiddleware{logContainer: logContainer}
}

func (m *LoggingMiddleware) WithLogging(next http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{0, 0}
		lw := loggingResponseWriter{responseWriter: w, data: responseData}
		next.ServeHTTP(&lw, r)
		uriAttr := slog.Attr{Key: "RequestURI", Value: slog.StringValue(r.RequestURI)}
		methodAttr := slog.Attr{Key: "Method", Value: slog.StringValue(r.Method)}
		durationAttr := slog.Attr{Key: "Duration", Value: slog.StringValue(time.Since(start).String())}

		statusCodeAttr := slog.Attr{Key: "StatusCode", Value: slog.IntValue(lw.data.statusCode)}
		sizeAttr := slog.Attr{Key: "ResponseSize", Value: slog.IntValue(lw.data.size)}
		m.logContainer.Logger.Log(r.Context(), slog.LevelInfo, "HTTP Request", uriAttr, methodAttr, durationAttr)
		m.logContainer.Logger.Log(r.Context(), slog.LevelInfo, "HTTP Response", statusCodeAttr, sizeAttr)
	}

	return http.HandlerFunc(logFn)
}
