package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/naspinall/Hive/pkg/models"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")

		if bearerToken == "" {
			http.Error(w, "Token Required", http.StatusUnauthorized)
			return
		}

		splitStrings := strings.SplitAfter(bearerToken, "Bearer ")

		if len(splitStrings) < 1 {
			http.Error(w, "Token Required", http.StatusUnauthorized)
			return
		}
		tokenFromHeader := splitStrings[1]
		token, err := jwt.ParseWithClaims(
			tokenFromHeader,
			&models.UserClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte("JWTSecretReallySecret"), nil
			},
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*models.UserClaims)
		if !ok {
			http.Error(w, "Bad Claims", http.StatusUnauthorized)
			return
		}

		if claims.ExpiresAt < time.Now().UTC().Unix() {
			http.Error(w, "Token Expired", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
