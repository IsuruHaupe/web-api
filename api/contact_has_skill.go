package api

import (
	"net/http"

	auth "github.com/IsuruHaupe/web-api/auth/token"
	db "github.com/IsuruHaupe/web-api/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createContactHasSkillRequest struct {
	ContactID int32 `json:"contact_id" binding:"required"`
	SkillID   int32 `json:"skill_id" binding:"required"`
}

// createSkillToContact godoc
// @Security bearerAuth
// @Summary Create a skill for a contact
// @Description This function is used to bind a skill to a contact.
// @Tags Bind Skill To Contact
// @Accept json
// @Produce json
// @Param createContactHasSkill body api.createContactHasSkillRequest true "Create Skill To Contact"
// @Success 200 {object} db.ContactHasSkill
// @Router /add-skill [post]
func (server *Server) createSkillToContact(ctx *gin.Context) {
	var req createContactHasSkillRequest
	// We verify that the JSON is correct, i.e : all fields are present.
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Retrieve the username in the authorization payload.
	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Payload)
	// In case of no error, we add the contact in the database.
	args := db.CreateContactHasSkillParams{
		Owner:     authPayload.Username,
		ContactID: req.ContactID,
		SkillID:   req.SkillID,
	}

	contact, err := server.database.CreateContactHasSkill(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, contact)
}
