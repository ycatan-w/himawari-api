package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ycatan-w/himawari-api/internal/api/internal/utils"
	"github.com/ycatan-w/himawari-api/internal/db"
)

type userContextKey struct{}

func GetUserID(r *http.Request) (int, bool) {
	userId, ok := r.Context().Value(userContextKey{}).(int)
	return userId, ok
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.RespondJSON(w, http.StatusUnauthorized, utils.CreateResponseError(
				"requires_token",
				[]utils.APIError{utils.BuildApiError("", "REQUIRED_TOKEN", "A Token is required.")},
			))

			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		userId, err := db.ValidateToken(token)
		if userId == -1 || err != nil {
			utils.RespondJSON(w, http.StatusUnauthorized, utils.CreateResponseError(
				"invalid_token",
				[]utils.APIError{utils.BuildApiError("", "INVALID_TOKEN", "The Token is invalid.")},
			))

			return
		}

		ctx := context.WithValue(r.Context(), userContextKey{}, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func ValidateRouteMethod(w http.ResponseWriter, r *http.Request, allowed map[string]bool, methods []string) bool {
	if !allowed[r.Method] {
		utils.RespondJSON(w, http.StatusMethodNotAllowed, utils.CreateResponseError(
			"method_not_allowed",
			[]utils.APIError{utils.BuildApiError("", "METHOD_NOT_ALLOWED", fmt.Sprintf("Expected Method %v but got %s", methods, r.Method))},
		))
		return false
	}
	return true
}
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
