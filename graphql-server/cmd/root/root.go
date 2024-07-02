package root

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/maximrub/thesis-diseases-risk-factors-server/internal/auth"
	"github.com/maximrub/thesis-diseases-risk-factors-server/internal/dal"
	"github.com/maximrub/thesis-diseases-risk-factors-server/internal/generated"
	"github.com/maximrub/thesis-diseases-risk-factors-server/internal/graph"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
)

const (
	defaultPort = 80
)

var (
	port int
)

var rootCmd = &cobra.Command{
	Use:   "diseases-risk-factors",
	Short: "Diseases Risk Factors",
	Run: func(cmd *cobra.Command, args []string) {
		run(cmd.Context())
	},
}

func init() {
	rootCmd.PersistentFlags().
		IntVar(&port, "port", defaultPort, "Listen port")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(fmt.Errorf("failed to execute: %w", err))
	}
}

func run(ctx context.Context) {
	dal, err := dal.NewDal(ctx)
	if err != nil {
		log.WithError(err).Fatal("error init dal")
	}

	defer func() {
		if err = dal.Close(context.TODO()); err != nil {
			log.WithError(err).Fatal("error closing dal")
		}
	}()

	resolver := &graph.Resolver{
		DAL: dal,
	}

	router := chi.NewRouter()

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	authenticatedPattern := "/query"
	router.Handle("/", playground.Handler("GraphQL Playground", authenticatedPattern))
	router.Mount(authenticatedPattern, authenticatedRouter(srv))

	log.Infof("connect to http://localhost:%d/ for GraphQL", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		log.WithError(err).Fatalf("error init server on [:%d]", port)
	}
}

// A completely separate router for authenticated routes
func authenticatedRouter(srv *handler.Server) http.Handler {
	r := chi.NewRouter()

	// add the middleware for checking the JWT token
	r.Use(auth.EnsureValidToken())

	// add the middlewares for checking the scopes in the JWT claims
	r.Use(auth.Authorizer())

	r.Handle("/", srv)
	return r
}
