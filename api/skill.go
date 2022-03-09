package api

import (
	"database/sql"
	"net/http"

	db "github.com/IsuruHaupe/web-api/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createSkillRequest struct {
	SkillName  string `json:"skill_name" binding:"required"`
	SkillLevel string `json:"skill_level" binding:"required"`
}

func (server *Server) createSkill(ctx *gin.Context) {
	var req createSkillRequest
	// We verify that the JSON is correct, i.e : all fields are present.
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// In case of no error, we add the contact in the database.
	args := db.CreateSkillParams{
		SkillName:  req.SkillName,
		SkillLevel: req.SkillLevel,
	}

	contact, err := server.database.CreateSkill(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, contact)
}

type getSkillRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getSkill(ctx *gin.Context) {
	var req getSkillRequest
	// We verify that the JSON is correct, i.e : all fields are present.
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	skill, err := server.database.GetSkill(ctx, req.ID)
	if err != nil {
		// Check if we have no skill with that ID.
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, skill)
}

type listSkillsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=10"`
}

func (server *Server) listSkills(ctx *gin.Context) {
	var req listSkillsRequest
	// We verify that the JSON is correct, i.e : all fields are present.
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.ListSkillsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	skills, err := server.database.ListSkills(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, skills)
}

type deleteSkillRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteSkill(ctx *gin.Context) {
	var req deleteSkillRequest
	// We verify that the JSON is correct, i.e : all fields are present.
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.database.DeleteSkill(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Successfully deleted skill.")
}

type updateSkillRequest struct {
	ID         int64  `json:"id" binding:"required,min=1"`
	SkillName  string `json:"skill_name,omitempty"`
	SkillLevel string `json:"skill_level,omitempty"`
}

func skillPreviousValues(ctx *gin.Context, req *updateSkillRequest, server *Server) {
	if req.SkillName == "" {
		skillName, err := server.database.GetSkillName(ctx, req.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		req.SkillName = skillName
	}

	if req.SkillLevel == "" {
		skillLevel, err := server.database.GetSkillLevel(ctx, req.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		req.SkillLevel = skillLevel
	}

}

func (server *Server) updateSkill(ctx *gin.Context) {
	var req updateSkillRequest
	// We verify that the JSON is correct, i.e : all fields are present.
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// For each filed we check if it nil.
	// If it is, then we retrieve the previous value for it.
	skillPreviousValues(ctx, &req, server)

	// In case of no error, we update the skill in the database.
	args := db.UpdateSkillParams{
		ID:         req.ID,
		SkillName:  req.SkillName,
		SkillLevel: req.SkillLevel,
	}

	contact, err := server.database.UpdateSkill(ctx, args)
	if err != nil {
		// Check if we have no contact with that ID.
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, contact)
}
