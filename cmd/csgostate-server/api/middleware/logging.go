package middleware

import (
	"github.com/gorilla/mux"
	"github.com/roessland/csgostate/cmd/csgostate-server/server"
	"net/http"
)

func NewRequestResponseLoggingMiddleware(app *server.App) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			route := mux.CurrentRoute(r)
			var pathTemplate string
			if route != nil {
				pathTemplate, _ = route.GetPathTemplate()
			}
			app.Log.Infow("http.request",
				"method", r.Method,
				"route", pathTemplate,
				"uri", r.RequestURI)

			srw := &statusResponseWriter{ResponseWriter: w}

			next.ServeHTTP(srw, r)

			app.Log.Infow("http.response",
				"method", r.Method,
				"route", pathTemplate,
				"uri", r.RequestURI,
				"statusCode", srw.statusCode)
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
