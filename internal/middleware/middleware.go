package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/supabase-community/supabase-go"
)

// ContextKey is the type for context keys used by this package.
type ContextKey string

// UserIDKey is the context key for the authenticated user's ID.
const UserIDKey ContextKey = "user_id"

// GetUserID retrieves the authenticated user ID stored by RequireAuth.
func GetUserID(ctx context.Context) string {
	if id, ok := ctx.Value(UserIDKey).(string); ok {
		return id
	}
	return ""
}

// TokenVerifier abstracts session token validation, making RequireAuth testable.
type TokenVerifier interface {
	VerifyToken(token string) (string, error)
}

// SupabaseTokenVerifier wraps *supabase.Client to implement TokenVerifier.
type SupabaseTokenVerifier struct {
	Client *supabase.Client
}

func (s *SupabaseTokenVerifier) VerifyToken(token string) (string, error) {
	user, err := s.Client.Auth.WithToken(token).GetUser()
	if err != nil {
		return "", err
	}
	return user.ID.String(), nil
}

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

// WriteHeader captures the status code and calls the underlying WriteHeader
func (rw *responseWriter) WriteHeader(statusCode int) {
	if !rw.written {
		rw.statusCode = statusCode
		rw.written = true
		rw.ResponseWriter.WriteHeader(statusCode)
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.statusCode = http.StatusOK
		rw.written = true
	}
	return rw.ResponseWriter.Write(b)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the ResponseWriter to capture status code
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // default to 200 if not set
			written:        false,
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		log.Printf(
			"%s %s from %s - took %v - status code - %d",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			duration,
			rw.statusCode,
		)
	})
}

// SecurityHeaders adds security-related HTTP headers
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		next.ServeHTTP(w, r)
	})
}

func RequireAuth(verifier TokenVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session")
			if err != nil {
				http.Redirect(w, r, "/adminlogin", http.StatusSeeOther)
				return
			}
			userID, err := verifier.VerifyToken(cookie.Value)
			if err != nil {
				http.Redirect(w, r, "/adminlogin", http.StatusSeeOther)
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
