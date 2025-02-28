//go:build !solution

package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type contextKey string

const userKey contextKey = "user"

var ErrInvalidToken = errors.New("invalid token")

type User struct {
	Name  string
	Email string
}

type TokenChecker interface {
	CheckToken(context context.Context, token string) (*User, error)
}

func ContextUser(context context.Context) (*User, bool) {
	user, ok := context.Value(userKey).(*User)
	return user, ok
}

func CheckAuth(checker TokenChecker) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := getTokenFromRequest(r)
			if err != nil {
				handleError(w, err)
				return
			}

			user, err := authenticateToken(r.Context(), checker, token)
			if err != nil {
				handleError(w, err)
				return
			}

			contextParameter := context.WithValue(r.Context(), userKey, user)
			next.ServeHTTP(w, r.WithContext(contextParameter))
		})
	}
}

func getTokenFromRequest(request *http.Request) (string, error) {
	token := request.Header.Get("Authorization")
	if token == "" {
		return "", errors.New("authorization token is missing")
	}

	if !strings.HasPrefix(token, "Bearer ") {
		return "", errors.New("invalid token format")
	}
	return token[len("Bearer "):], nil
}

func authenticateToken(context context.Context, checker TokenChecker, token string) (*User, error) {
	user, err := checker.CheckToken(context, token)
	if err != nil {
		if errors.Is(err, ErrInvalidToken) {
			return nil, errors.New("invalid token")
		}
		return nil, errors.New("internal server error")
	}
	return user, nil
}

func handleError(writer http.ResponseWriter, err error) {
	switch {
	case strings.Contains(err.Error(), "authorization token is missing"):
		http.Error(writer, "Authorization token is missing", http.StatusUnauthorized)
	case strings.Contains(err.Error(), "invalid token format"):
		http.Error(writer, "Invalid token format", http.StatusUnauthorized)
	case strings.Contains(err.Error(), "invalid token"):
		http.Error(writer, "Invalid token", http.StatusUnauthorized)
	case strings.Contains(err.Error(), "internal server error"):
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
	default:
		http.Error(writer, "Unexpected error", http.StatusInternalServerError)
	}
}
