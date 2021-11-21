package middleware

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/roessland/csgostate/internal/server"
	"net/http"
	"strconv"
	"time"
)

var httpRequestHistogram = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:      "http_request_duration_seconds",
		Help:      "Duration of HTTP requests",
	}, []string{
		"app", "env", "handler", "method", "code",
	},
)

func NewRequestResponseLoggingMiddleware(app *server.App) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t0 := time.Now()

			// Extract mux path matched
			route := mux.CurrentRoute(r)
			var pathTemplate string
			if route != nil {
				pathTemplate, _ = route.GetPathTemplate()
			}

			sess, _ := app.SessionStore.New(r)
			steamID := sess.SteamID()
			nickName := sess.NickName()

			app.Log.Infow("http.request",
				"method", r.Method,
				"route", pathTemplate,
				"uri", r.RequestURI,
				"steamid", steamID,
				"nickname", nickName,
				"request_id", RequestID(r),
			)

			srw := &statusResponseWriter{ResponseWriter: w}
			next.ServeHTTP(srw, r)

			duration := time.Since(t0)

			httpRequestHistogram.
				WithLabelValues(
					app.Config.AppName,
					app.Config.Env,
					pathTemplate,
					r.Method,
					strconv.Itoa(srw.statusCode)).
				Observe(duration.Seconds())

			app.Log.Infow("http.response",
				"method", r.Method,
				"route", pathTemplate,
				"uri", r.RequestURI,
				"statusCode", srw.statusCode,
				"request_id", RequestID(r),
				"duration_ms", duration.Milliseconds(),
			)
		})
	}
}

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (srw *statusResponseWriter) WriteHeader(statusCode int) {
	srw.statusCode = statusCode
	srw.ResponseWriter.WriteHeader(statusCode)
}

func (srw *statusResponseWriter) Write(data []byte) (int, error) {
	if srw.statusCode == 0 {
		srw.statusCode = http.StatusOK
	}
	return srw.ResponseWriter.Write(data)
}
