package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/Ruturaj-7802/taskflow/internal/service"
)

type contextKey string

const UserIDKey contextKey = "user_id"
const UserClaimsKey contextKey = "user_claims"

func Auth(authSvc *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(header, "Bearer ")

			claims, err := authSvc.ValidateToken(tokenStr)
			if err != nil {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			// ✅ claims.UserID is uuid.UUID
			userID := claims.UserID
			if userID == uuid.Nil {
				http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				return
			}

			// store both claims and userID in context
			ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
			ctx = context.WithValue(ctx, UserIDKey, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}