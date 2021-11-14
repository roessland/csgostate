package middleware

import (
	"context"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/roessland/csgostate/cmd/csgostate-server/server"
	"net/http"
)

type contextKey string

const contextKeyRequestID contextKey = "request-id"

func WithRequestID(r *http.Request, requestID string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), contextKeyRequestID, requestID))
}

func RequestID(r *http.Request) string {
	return r.Context().Value(contextKeyRequestID).(string)
}

func NewRequestIDMiddleware(app *server.App) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := uuid.New().String()
			next.ServeHTTP(w, WithRequestID(r, requestID))
		})
	}
}
