package auth

import (
	"context"
	"errors"
	"github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/maximrub/thesis-diseases-risk-factors-server/internal/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"os"
	"time"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var scopeCtxKey = &contextKey{
	name: "scope",
}

type contextKey struct {
	name string
}

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	Scope       string   `json:"scope"`
	Permissions []string `json:"permissions"`
}

// Validate does nothing for this example, but we need
// it to satisfy validator.CustomClaims interface.
func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

// NotAuthorizedError is a custom error to be returned when the authorization fails
var NotAuthorizedError = errors.New("not authorized")

// EnsureValidToken is a middleware that will check the validity of our JWT.
func EnsureValidToken() func(next http.Handler) http.Handler {
	issuerURL, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/")
	if err != nil {
		log.WithError(err).Fatal("Failed to parse the issuer url")
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{os.Getenv("AUTH0_AUDIENCE")},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		log.WithError(err).Fatal("Failed to set up the jwt validator")
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		log.WithError(err).Error("Encountered error while validating JWT")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Failed to validate JWT."}`))
	}

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	return func(next http.Handler) http.Handler {
		return middleware.CheckJWT(next)
	}
}

func Authorizer() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
			if !ok {
				http.Error(w, "failed to get validated claims", http.StatusInternalServerError)
				return
			}

			customClaims, ok := claims.CustomClaims.(*CustomClaims)
			if !ok {
				http.Error(w, "could not cast custom claims to specific type", http.StatusInternalServerError)
			}

			// put it in context
			ctx := context.WithValue(r.Context(), scopeCtxKey, customClaims)

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func HasPermission(ctx context.Context, requiredPermission string) bool {
	jwtClaims := ctx.Value(scopeCtxKey).(*CustomClaims)
	return utils.Contains(jwtClaims.Permissions, requiredPermission)
}
