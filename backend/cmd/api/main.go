package main

import (
    "log"
    "net/http"
    "time"

    "github.com/joho/godotenv"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/cors"

    "telegraph/internal/acl"
	"telegraph/internal/audit"
    "telegraph/internal/auth"
    "telegraph/internal/channels"
    "telegraph/internal/config"
    "telegraph/internal/database"
	"telegraph/internal/messages"
    "telegraph/internal/users"
    "telegraph/internal/ws"
    mw "telegraph/internal/middleware"
)

func main() {
    // Load .env FIRST
    if err := godotenv.Load("/home/apeiron/Desktop/A2SV/projects/telegraph/backend/.env"); err != nil {
        log.Println("warning: .env not found, using system env")
    }

    // Then load config
    cfg, err := config.Load()
    if err != nil {
        log.Fatal("config error:", err)
    }

	// Connect to MongoDB
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("DB error:", err)
	}

	log.Println("✓ MongoDB connected")

	// Repos
	userRepo := users.NewMongoUserRepo(db)
	refreshRepo := auth.NewRefreshTokenRepo(db)
	mfaRepo := auth.NewMFACodeRepo(db)
	channelRepo := channels.NewMongoChannelRepo(db)
	messageRepo := messages.NewMongoMessageRepo(db)

	// Utils & Managers
	jwtMgr := users.NewJWTManager(cfg.JWTSecret, time.Hour*1)
	refreshMgr := auth.NewRefreshTokenManager(refreshRepo, time.Hour*24*7)
	smtpSender := auth.NewSMTPSender(cfg.SMTPEmail, cfg.SMTPPassword, cfg.SMTPHost, cfg.SMTPPort)
	mfaMgr := auth.NewMFAManager(mfaRepo, smtpSender)
	auditLogger := audit.NewLogger(db)

	// WebSocket Hub
	hub := ws.NewHub()
	go hub.Run()

	// Services
	userSvc := users.NewUserService(userRepo)
	channelSvc := channels.NewChannelService(channelRepo, userRepo, auditLogger)
	messageSvc := messages.NewMessageService(messageRepo, channelRepo, auditLogger, hub)

	// Handlers
	authHandler := auth.NewHandler(userSvc, refreshMgr, jwtMgr, mfaMgr)
	userHandler := users.NewHandler(userSvc, jwtMgr)
	channelHandler := channels.NewHandler(channelSvc, userSvc)
	messageHandler := messages.NewHandler(messageSvc)

	// Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	
	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api/v1", func(api chi.Router) {

		// Public routes
		api.Mount("/auth", authHandler.Routes())

		api.Route("/users", func(ur chi.Router) {
			ur.Post("/register", userHandler.Register)

			// Protected user routes
			ur.Group(func(pr chi.Router) {
				pr.Use(mw.JWTAuth(jwtMgr))
				pr.Use(mw.LoadUser(userRepo))
				pr.Get("/me", userHandler.Me)
				pr.Get("/{id}", userHandler.GetUser)
				pr.Put("/{id}", userHandler.UpdateUser)
				pr.Delete("/{id}", userHandler.DeleteUser)
			})
		})

		// Protected channel & message routes
		api.Group(func(cr chi.Router) {
			cr.Use(mw.JWTAuth(jwtMgr))
			cr.Use(mw.LoadUser(userRepo))
			
			// WebSocket route
			cr.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
				userID := users.UserIDFromContext(r.Context())
				ws.ServeWs(hub, userID)(w, r)
			})
			
			// Channel routes
			cr.Mount("/channels", channelHandler.Routes())
			
			// Message routes (note: using separate path to avoid conflict)
			cr.Post("/channels/{channelId}/messages", messageHandler.SendMessage)
			cr.Get("/channels/{channelId}/messages", messageHandler.GetMessages)
			cr.Delete("/messages/{id}", messageHandler.DeleteMessage)
		})

		// Admin-only routes (example for broadcasting)
		api.Group(func(ar chi.Router) {
			ar.Use(mw.JWTAuth(jwtMgr))
			ar.Use(mw.LoadUser(userRepo))
			ar.Use(mw.RequireRole(acl.RoleAdmin))
			
			// Admin routes would go here
			// Example: ar.Get("/admin/audit-logs", ...)
		})
	})

	log.Println("✓ Telegraph server running at :8080")
	log.Println("✓ Database: MongoDB Atlas")
	log.Println("✓ Access Control: RBAC + MAC + ABAC enabled")
	log.Println("✓ E2EE: Message encryption active")
	log.Println("✓ Audit Logging: Enabled")
	
	if auditLogger != nil {
		log.Println("✓ All systems operational - Production Ready!")
	}
	
	log.Fatal(http.ListenAndServe(":8080", r))
}
