package logger

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// New - настройка собственного логирования для middleware
func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// параметр, который выводится с каждой строкой логов
		log = log.With(
			slog.String("component", "middleware/logger"),
		)

		log.Info("logger middleware enabled")

		// выводится при каждом запросе (в цепочке handlers)
		fn := func(w http.ResponseWriter, r *http.Request) {

			// (начало) до обработки запроса
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)
			
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()

			// (конец) после окончательной обработки запроса
			defer func() {
				entry.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(t1).String()),
				)
			}()
			
			// (середина) передаем управление следующему handler в цепочке (may be следующему middleware)
			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}