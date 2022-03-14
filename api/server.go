package api

import (
	"fmt"

	auth "github.com/IsuruHaupe/web-api/auth/token"
	"github.com/IsuruHaupe/web-api/config"
	"github.com/IsuruHaupe/web-api/db/database"
	"github.com/IsuruHaupe/web-api/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// This struct is used to regroup the database connection, the gin router and the configuration.
type Server struct {
	config     config.Config
	database   database.Database
	tokenMaker auth.Maker
	router     *gin.Engine
}

// This function will create a new server and setup all routes.
func NewServer(config config.Config, database database.Database) (*Server, error) {
	tokenMaker, err := auth.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker : %w", err)
	}
	server := &Server{
		config:     config,
		database:   database,
		tokenMaker: tokenMaker,
	}
	server.setUpRouter()
	return server, nil
}

// This function will setup all routes.
func (server *Server) setUpRouter() {
	// Swagger 2.0 Meta Information.
	docs.SwaggerInfo.Title = "Web API."
	docs.SwaggerInfo.Description = "Web API for managing skills and contacts."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

	router := gin.Default()

	// Group routes that need authentification/authorization together.
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	// Add routes to the gin server.
	// Contacts routes.
	authRoutes.POST("/contacts", server.createContact)
	authRoutes.GET("/contacts/:id", server.getContact)
	authRoutes.GET("/contact-skills/:id", server.getContactSkills)
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
	// Binding skills and contacts route.
	authRoutes.POST("/add-skill", server.createSkillToContact)
	// Authentification routes.
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/tokens/renew_access", server.renewAccessToken)
	// Documentation routes, available at : http://localhost:8080/swagger/index.html.
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	server.router = router
}

// Start the server on the specified address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// Wrapper for returning error.
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
