package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/naspinall/Hive/pkg/models"
)

func WithFilter() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			f, err := models.NewFilterFromQueryString(r.URL.Query())
			if err != nil {
				http.Error(w, "Bad Filter Parameters Provided", http.StatusBadRequest)
				return
			}
			ctx := context.WithValue(r.Context(), models.FilterContextKey("Filter"), f)
			// Adding context to request of processing
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
