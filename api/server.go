package api

import (
	"fmt"

	auth "github.com/IsuruHaupe/web-api/auth/token"
	"github.com/IsuruHaupe/web-api/config"
	"github.com/IsuruHaupe/web-api/db/postgres"
	"github.com/gin-gonic/gin"
)

// This struct is used to regroup the database connection and the gin router.
type Server struct {
	config     config.Config
	database   postgres.Database
	tokenMaker auth.Maker
	router     *gin.Engine
}

// This function will create a new server and setup all routes.
func NewServer(config config.Config, database postgres.Database) (*Server, error) {
	tokenMaker, err := auth.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token make : %w", err)
	}
	server := &Server{
		config:     config,
		database:   database,
		tokenMaker: tokenMaker,
	}
	server.setUpRouter()
	return server, nil
}

func (server *Server) setUpRouter() {
	router := gin.Default()

	// Group routes that need authentification/authorization together.
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	// Add routes to the gin server.
	// Contacts routes.
	authRoutes.POST("/contacts", server.createContact)
	authRoutes.GET("/contacts/:id", server.getContact)
	authRoutes.GET("/contacts", server.listContacts)
	authRoutes.GET("/contacts-with-skill", server.getContactWithSkill)
	authRoutes.GET("/contacts-with-skill-and-level", server.getContactWithSkillAndLevel)
	authRoutes.DELETE("/contacts/:id", server.deleteContact)
	authRoutes.PATCH("/contacts", server.updateContact)
	// Skills routes.
	authRoutes.POST("/skills", server.createSkill)
	authRoutes.GET("/skills/:id", server.getSkill)
	authRoutes.GET("/skills", server.listSkills)
	authRoutes.DELETE("/skills/:id", server.deleteSkill)
	authRoutes.PATCH("/skills", server.updateSkill)
	// Binding skills and contacts route
	authRoutes.POST("/add-skill", server.createSkillToContact)
	// Authentification routes
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	server.router = router
}

// Start the server on the specified address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
