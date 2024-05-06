package middleware

import (
	"context"
	"fmt"
	"net/http"
	"ps-cats-social/pkg/httphelper"
	"ps-cats-social/pkg/httphelper/response"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/exp/slog"
)

var jwtKey = []byte("your_secret_key")

func JWTAuthMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		jwtToken, err := extractJWTTokenFromHeader(r)
		if err != nil {
			slog.Error("Failed to extract JWT token from header", "error", err)
			writeUnauthorized(rw)
			return
		}

		claims, err := parseJWTToClaims(jwtToken)
		if err != nil {
			slog.Error("Failed to parse JWT token", "error", err)
			writeUnauthorized(rw)
			return
		}

		email, emailOk := claims["email"].(string)
		userId, uidOk := claims["user_id"].(float64)
		if !emailOk || !uidOk {
			slog.Error("Invalid claims")
			writeUnauthorized(rw)
			return
		}

		user, err := constructUserInfo(email, userId)
		if err != nil {
			slog.Error("Failed to construct user info", "error", err)
			writeUnauthorized(rw)
			return
		}

		r2 := r.WithContext(context.WithValue(r.Context(), "user_info", user))
		slog.Debug("AUTHORIZED", "email", r2.Context().Value("email"))
		slog.Debug("AUTHORIZED", "user_id", r2.Context().Value("user_id"))

		fn(rw, r2)
	}
}

func writeUnauthorized(rw http.ResponseWriter) {
	httphelper.WriteJSON(
		rw, http.StatusUnauthorized,
		response.WebResponse{
			Status:  http.StatusUnauthorized,
			Message: http.StatusText(http.StatusUnauthorized),
		},
	)
}

func extractJWTTokenFromHeader(r *http.Request) (string, error) {
	authToken := r.Header.Get("Authorization")
	if authToken == "" {
		return "", fmt.Errorf("missing auth token")
	}

	return authToken[len("Bearer "):], nil
}

type Claims struct {
	Email  string `json:"email"`
	UserId int64  `json:"user_id"`
	jwt.Claims
}

func GenerateJWT(email string, userId int64) (string, error) {
	// Create token
	claims := Claims{
		Email:  email,
		UserId: userId,
		Claims: jwt.MapClaims{
			"exp": time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and return it
	return token.SignedString(jwtKey)
}

func constructUserInfo(email string, userId float64) (map[string]interface{}, error) {
	return map[string]interface{}{
		"email":   email,
		"user_id": int64(userId),
	}, nil
}

func parseJWTToClaims(jwtToken string) (jwt.MapClaims, error) {
	token, _, err := jwt.NewParser().ParseUnverified(jwtToken, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	// no need to verify 'token' signature since it already validated in authz kong plugin, just parse the token

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid jwt token")
	}

	return claims, nil
}
