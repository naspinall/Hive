package middleware

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/naspinall/Hive/pkg/models"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Token Required", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(
			cookie.Value,
			&models.UserClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte("JWTSecretReallySecret"), nil
			},
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		_, ok := token.Claims.(*models.UserClaims)
		if !ok {
			http.Error(w, "Bad Claims", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
