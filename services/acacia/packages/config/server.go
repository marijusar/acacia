package config

import (
	"acacia/packages/api"
	"acacia/packages/auth"
	"acacia/packages/crypto"
	"acacia/packages/db"
	"acacia/packages/llm"
	"acacia/packages/routes"
	"acacia/packages/services"
	"acacia/packages/storage"
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

type Database struct {
	Queries *db.Queries
	Conn    *sql.DB
}

type Server struct {
	httpServer *http.Server
	db         *Database
	logger     *logrus.Logger
}

func NewServer(d *Database, l *logrus.Logger, env *Environment) *Server {
	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(
		env.JWTSecret,
		15*time.Minute,  // Access token: 15 minutes
		30*24*time.Hour, // Refresh token: 30 days
	)

	// Initialize encryption service
	encryptionService, err := crypto.NewEncryptionService(env.EncryptionKey)
	if err != nil {
		l.WithError(err).Fatal("Failed to initialize encryption service")
	}

	// Initialize authentication middleware
	authMiddleware := auth.NewAuthMiddleware(jwtManager, d.Queries, l)

	// Initialize authorization middleware (single instance)
	authzMiddleware := auth.NewAuthorizationMiddleware(d.Queries, l)

	// Initialize LLM provider registry
	providerRegistry := llm.NewProviderRegistry()

	// Initialize conversation service
	conversationService := services.NewConversationService(
		d.Queries,
		providerRegistry,
		encryptionService,
		l,
	)

	// Initialize S3 storage
	s3Storage, err := storage.NewS3Storage(storage.S3Config{
		Bucket:          env.AWSS3Bucket,
		Region:          env.AWSRegion,
		AccessKeyID:     env.AWSAccessKeyID,
		SecretAccessKey: env.AWSSecretKey,
		Endpoint:        env.AWSEndpoint,
	}, l)
	if err != nil {
		l.WithError(err).Fatal("Failed to initialize S3 storage")
	}

	issuesController := api.NewIssuesController(d.Queries, l, s3Storage)
	projectsController := api.NewProjectsController(d.Queries, l)
	projectColumnsController := api.NewProjectStatusColumnsController(d.Queries, l, d.Conn)
	usersController := api.NewUsersController(d.Queries, l, jwtManager)
	teamsController := api.NewTeamsController(d.Queries, l)
	teamLLMAPIKeysController := api.NewTeamLLMAPIKeysController(d.Queries, l, encryptionService)
	conversationsController := api.NewConversationsController(d.Queries, l, conversationService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	authMiddlewares := chi.Middlewares{authMiddleware.Handle}

	r.Mount("/issues", routes.IssuesRoutes(issuesController, authMiddlewares, authzMiddleware))
	r.Mount("/projects", routes.ProjectsRoutes(projectsController, authMiddlewares, authzMiddleware))
	r.Mount("/project-columns", routes.ProjectStatusColumnsRoutes(projectColumnsController, authMiddlewares, authzMiddleware))
	r.Mount("/users", routes.UsersRoutes(usersController, authMiddlewares))
	r.Mount("/teams", routes.TeamsRoutes(teamsController, teamLLMAPIKeysController, authMiddlewares, authzMiddleware))
	r.Mount("/conversations", routes.ConversationsRoutes(conversationsController, authMiddlewares, authzMiddleware))

	httpServer := &http.Server{
		Handler: r,
	}

	server := &Server{
		httpServer: httpServer,
		db:         d,
		logger:     l,
	}

	return server
}

func (s *Server) ListenAndServe(port string) {
	s.httpServer.Addr = ":" + port
	s.logger.WithField("port", port).Info("Starting server")
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.WithError(err).Error("Server failed to start")
	}
}

func (s *Server) Close() {
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*5000)
	s.httpServer.Shutdown(ctx)
}
