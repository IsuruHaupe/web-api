package api

import (
	"github.com/IsuruHaupe/web-api/postgres"
	"github.com/gin-gonic/gin"
)

// This struct is used to regroup the database connection and the gin router.
type Server struct {
	database *postgres.PostgresDatabase
	router   *gin.Engine
}

// This function will create a new server and setup all routes.
func NewServer(database *postgres.PostgresDatabase) *Server {
	server := &Server{
		database: database,
	}
	router := gin.Default()

	// Add routes to the gin server.
	// Contacts routes.
	router.POST("/contacts", server.createContact)
	router.GET("/contacts/:id", server.getContact)
	router.GET("/contacts", server.listContacts)
	router.GET("/contacts-with-skill", server.getContactWithSkill)
	router.GET("/contacts-with-skill-and-level", server.getContactWithSkillAndLevel)
	router.DELETE("/contacts/:id", server.deleteContact)
	router.PATCH("/contacts", server.updateContact)
	// Skills routes.
	router.POST("/skills", server.createSkill)
	router.GET("/skills/:id", server.getSkill)
	router.GET("/skills", server.listSkills)
	router.DELETE("/skills/:id", server.deleteSkill)
	router.PATCH("/skills", server.updateSkill)
	// Binding skills and contacts route
	router.POST("/add-skill", server.createSkillToContact)

	server.router = router
	return server
}

// Start the server on the specified address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
