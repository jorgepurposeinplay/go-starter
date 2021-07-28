package logs

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	chimiddleware "github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewStructuredLogger(logger *zap.Logger) func(next http.Handler) http.Handler {
	return chimiddleware.RequestLogger(&structuredLogger{logger})
}

type structuredLogger struct {
	Logger *zap.Logger
}

func (l *structuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &structuredLoggerEntry{Logger: l.Logger}

	fields := []zapcore.Field{zap.String("ts", time.Now().UTC().Format(time.RFC1123))}

	fields = append(fields, []zapcore.Field{
		zap.String("component", "api"),
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
	}...)

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		fields = append(fields, zap.String("req.id", reqID))
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	fields = append(fields, []zapcore.Field{
		zap.String("http.scheme", scheme),
		zap.String("http.proto", r.Proto),
		zap.String("http.method", r.Method),
		zap.String("remote_addr", r.RemoteAddr),
		zap.String("user_agent", r.UserAgent()),
		zap.String("uri", fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)),
	}...)

	entry.Logger = l.Logger.With(fields...)

	entry.Logger.Info("request started")

	return entry
}

type structuredLoggerEntry struct {
	Logger *zap.Logger
}

func (l *structuredLoggerEntry) Write(status, bytes int, elapsed time.Duration) {
	l.Logger = l.Logger.With(
		zap.Int("res.status", status),
		zap.Int("res.bytes_length", bytes),
		zap.Float64("res.elapsed_ms", float64(elapsed.Nanoseconds())/1000000.0))

	l.Logger.Info("request complete")
}

func (l *structuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.With(
		zap.String("stack", string(stack)),
		zap.String("panic", fmt.Sprintf("%+v", v)),
	)
}

func (l *structuredLoggerEntry) WithError(err error) *zap.Logger {
	l.Logger = l.Logger.With(
		zap.Error(err),
	)
	return l.Logger
}

func GetLogEntry(r *http.Request) *structuredLoggerEntry {
	entry, _ := chimiddleware.GetLogEntry(r).(*structuredLoggerEntry)
	return entry
}
