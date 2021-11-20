package middleware

import (
	"github.com/gorilla/mux"
	"github.com/roessland/csgostate/internal/server"
	"net/http"
	"time"
)

func NewRequestResponseLoggingMiddleware(app *server.App) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract mux path matched
			route := mux.CurrentRoute(r)
			var pathTemplate string
			if route != nil {
				pathTemplate, _ = route.GetPathTemplate()
			}

			// Extract SteamID from Session
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

			t0 := time.Now()

			srw := &statusResponseWriter{ResponseWriter: w}

			next.ServeHTTP(srw, r)

			millis := time.Since(t0).Milliseconds()

			app.Log.Infow("http.response",
				"method", r.Method,
				"route", pathTemplate,
				"uri", r.RequestURI,
				"statusCode", srw.statusCode,
				"request_id", RequestID(r),
				"duration_ms", millis,
			)
		})
	}
}

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (srw *statusResponseWriter) WriteHeader(statusCode int) () {
	srw.statusCode = statusCode
	srw.ResponseWriter.WriteHeader(statusCode)
}

func (srw *statusResponseWriter) Write(data []byte) (int, error) {
	if srw.statusCode == 0 {
		srw.statusCode = http.StatusOK
	}
	return srw.ResponseWriter.Write(data)
}
