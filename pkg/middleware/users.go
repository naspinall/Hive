package middleware

import (
	"net/http"
	"regexp"

	"github.com/naspinall/Hive/pkg/models"

	"github.com/gorilla/mux"
)

type UsersMiddleware struct {
	us models.UserService
}

// Gets JWT from bearer token header.
var bearerTokenRegex = regexp.MustCompile(`Bearer ([A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_.+/=]*)`)

func NewUsersMiddleware(us models.UserService) *UsersMiddleware {
	return &UsersMiddleware{
		us: us,
	}
}

func (um *UsersMiddleware) JWTAuth() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headerToken := r.Header.Get("Authorization")
			user := &models.User{Token: headerToken}
			ctx, err := um.us.AcceptToken(user, r.Context())

			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// Adding context to request of processing
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
