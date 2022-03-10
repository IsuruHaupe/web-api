package api

import (
	"database/sql"
	"net/http"

	db "github.com/IsuruHaupe/web-api/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createContactRequest struct {
	Owner       string `json:"owner" binding:"required"`
	Firstname   string `json:"firstname" binding:"required"`
	Lastname    string `json:"lastname" binding:"required"`
	Fullname    string `json:"fullname" binding:"required"`
	HomeAddress string `json:"home_address" binding:"required"`
	Email       string `json:"email" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

func (server *Server) createContact(ctx *gin.Context) {
	var req createContactRequest
	// We verify that the JSON is correct, i.e : all fields are present.
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// In case of no error, we add the contact in the database.
	args := db.CreateContactParams{
		Owner:       req.Owner,
		Firstname:   req.Firstname,
		Lastname:    req.Lastname,
		Fullname:    req.Fullname,
		HomeAddress: req.HomeAddress,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
	}

	contact, err := server.database.CreateContact(ctx, args)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, contact)
}

type getContactRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getContact(ctx *gin.Context) {
	var req getContactRequest
	// We verify that the JSON is correct, i.e : all fields are present.
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	contact, err := server.database.GetContact(ctx, req.ID)
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

type listContactsRequest struct {
	// A SUPPRIMER
	Owner    string `json:"owner" binding:"required"`
	PageID   int32  `form:"page_id" binding:"required,min=1"`
	PageSize int32  `form:"page_size" binding:"required,min=1,max=10"`
}

func (server *Server) listContacts(ctx *gin.Context) {
	var req listContactsRequest
	// We verify that the JSON is correct, i.e : all fields are present.
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// TODO : recuperer le username via authentification
	args := db.ListContactsParams{
		Owner:  req.Owner, // a modifier par un autre username
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	contacts, err := server.database.ListContacts(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, contacts)
}

type deleteContactRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteContact(ctx *gin.Context) {
	var req deleteContactRequest
	// We verify that the JSON is correct, i.e : all fields are present.
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.database.DeleteContact(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Successfully deleted contact.")
}

type updateContactRequest struct {
	ID          int64  `json:"id" binding:"required,min=1"`
	Firstname   string `json:"firstname,omitempty"`
	Lastname    string `json:"lastname,omitempty"`
	Fullname    string `json:"fullname,omitempty"`
	HomeAddress string `json:"home_address,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

func contactPreviousValues(ctx *gin.Context, req *updateContactRequest, server *Server) {
	exists, err := server.database.GetIfExistsContactID(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if !exists {
		return
	}

	if req.Firstname == "" {
		firstname, err := server.database.GetFirstname(ctx, req.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		req.Firstname = firstname
	}

	if req.Lastname == "" {
		lastname, err := server.database.GetLastname(ctx, req.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		req.Lastname = lastname
	}

	if req.Fullname == "" {
		fullname, err := server.database.GetFullname(ctx, req.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		req.Fullname = fullname
	}

	if req.HomeAddress == "" {
		homeAddress, err := server.database.GetHomeAddress(ctx, req.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		req.HomeAddress = homeAddress
	}

	if req.Email == "" {
		email, err := server.database.GetEmail(ctx, req.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		req.Email = email
	}

	if req.PhoneNumber == "" {
		phoneNumber, err := server.database.GetPhoneNumber(ctx, req.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		req.PhoneNumber = phoneNumber
	}
}

func (server *Server) updateContact(ctx *gin.Context) {
	var req updateContactRequest
	// We verify that the JSON is correct, i.e : all fields are present.
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// For each filed we check if it nil.
	// If it is, then we retrieve the previous value for it.
	contactPreviousValues(ctx, &req, server)

	// In case of no error, we upodate the contact in the database.
	args := db.UpdateContactParams{
		ID:          req.ID,
		Firstname:   req.Firstname,
		Lastname:    req.Lastname,
		Fullname:    req.Fullname,
		HomeAddress: req.HomeAddress,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
	}

	contact, err := server.database.UpdateContact(ctx, args)
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

type getContactWithSkillRequest struct {
	SkillName string `form:"skill_name" binding:"required"`
}

func (server *Server) getContactWithSkill(ctx *gin.Context) {
	var req getContactWithSkillRequest
	// We verify that the JSON is correct, i.e : all fields are present.
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	contacts, err := server.database.GetContactsWithSkill(ctx, req.SkillName)
	if err != nil {
		// Check if we have no contact with that ID.
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, contacts)
}

type getContactWithSkillAndLevelRequest struct {
	SkillName  string `form:"skill_name" binding:"required"`
	SkillLevel string `form:"skill_level" binding:"required"`
}

func (server *Server) getContactWithSkillAndLevel(ctx *gin.Context) {
	var req getContactWithSkillAndLevelRequest
	// We verify that the JSON is correct, i.e : all fields are present.
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.GetContactsWithSkillAndLevelParams{
		SkillName:  req.SkillName,
		SkillLevel: req.SkillLevel,
	}

	contacts, err := server.database.GetContactsWithSkillAndLevel(ctx, args)
	if err != nil {
		// Check if we have no contact with that ID.
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, contacts)
}
