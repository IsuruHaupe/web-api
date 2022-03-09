package api

import (
	"net/http"

	db "github.com/IsuruHaupe/web-api/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createContactHasSkillRequest struct {
	ContactID int32 `json:"contact_id" binding:"required"`
	SkillID   int32 `json:"skill_id" binding:"required"`
}

func (server *Server) createSkillToContact(ctx *gin.Context) {
	var req createContactHasSkillRequest
	// We verify that the JSON is correct, i.e : all fields are present.
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// In case of no error, we add the contact in the database.
	args := db.CreateContactHasSkillParams{
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
