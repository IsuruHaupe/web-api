package api

import (
	"database/sql"
	"errors"
	"net/http"

	auth "github.com/IsuruHaupe/web-api/auth/token"
	db "github.com/IsuruHaupe/web-api/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// Request holder when receiving a create contact request.
type createContactRequest struct {
	Firstname   string `json:"firstname" binding:"required"`
	Lastname    string `json:"lastname" binding:"required"`
	Fullname    string `json:"fullname" binding:"required"`
	HomeAddress string `json:"home_address" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

// createContact godoc
// @Security bearerAuth
// @Summary Create a contact
// @Description This function is used to create a contact for an user.
// @Tags Contact
// @Accept json
// @Produce json
// @Param contact body api.createContactRequest true "Create Contact"
// @Success 200 {object} db.Contact
// @Router /contacts [post]
func (server *Server) createContact(ctx *gin.Context) {
	var req createContactRequest
	// We verify that the JSON is correct, i.e : all fields are present.
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Retrieve the username in the authorization payload.
	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Payload)
	// In case of no error, we add the contact in the database.
	args := db.CreateContactParams{
		Owner:       authPayload.Username,
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

// Request holder for get contact request.
type getContactRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// getContact godoc
// @Security bearerAuth
// @Summary Get a contact
// @Tags Contact
// @Description This function is used to get a contact for an user.
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} db.Contact
// @Router /contacts/{id} [get]
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

	// Check for owernership.
	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Payload)
	if contact.Owner != authPayload.Username {
		err := errors.New("contact doesn't belong to the user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, contact)
}

// Request holder for listing contact request.
type listContactsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=10"`
}

// listContacts godoc
// @Security bearerAuth
// @Summary List contacts
// @Tags Contact
// @Description This function is used to list contacts for an user.
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param page_id query int true "page_id"
// @Param page_size query int true "page_size"
// @Success 200 {array} db.Contact
// @Router /contacts [get]
func (server *Server) listContacts(ctx *gin.Context) {
	var req listContactsRequest
	// We verify that the JSON is correct, i.e : all fields are present.
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Check for authentification payload
	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Payload)
	args := db.ListContactsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	contacts, err := server.database.ListContacts(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Check for owernership.
	for _, contact := range contacts {
		if contact.Owner != authPayload.Username {
			err := errors.New("contact doesn't belong to the user")
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, contacts)
}

// Request holder for deleting contact request.
type deleteContactRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// deleteContact godoc
// @Security bearerAuth
// @Summary Delete a contact
// @Tags Contact
// @Description This function is used to delete a contact for an user.
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param id path int true "id"
// @Success 200 {string} string "Successfully deleted contact."
// @Router /contacts/{id} [delete]
func (server *Server) deleteContact(ctx *gin.Context) {
	var req deleteContactRequest
	// We verify that the JSON is correct, i.e : all fields are present.
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get the contact to check ownership before deletion.
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

	// Check for owernership.
	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Payload)
	if contact.Owner != authPayload.Username {
		err := errors.New("contact doesn't belong to the user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = server.database.DeleteContact(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Successfully deleted contact.")
}

// Request holder for updating contact request.
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

// updateContact godoc
// @Security bearerAuth
// @Summary Update a contact
// @Tags Contact
// @Description This function is used to update a contact for an user.
// @Accept json
// @Produce json
// @Param contact body api.updateContactRequest true "Update Contact"
// @Success 200 {object} db.Contact
// @Router /contacts [patch]
func (server *Server) updateContact(ctx *gin.Context) {
	var req updateContactRequest
	// We verify that the JSON is correct, i.e : all fields are present.
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// Get the contact to check ownership before update.
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

	// Check for owernership.
	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Payload)
	if contact.Owner != authPayload.Username {
		err := errors.New("contact doesn't belong to the user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// For each field we check if it nil.
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

	contact, err = server.database.UpdateContact(ctx, args)
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

// Request holder for get contact with a specific skill request.
type getContactWithSkillRequest struct {
	SkillName string `form:"skill_name" binding:"required"`
}

// getContactWithSkill godoc
// @Security bearerAuth
// @Summary Get all the contacts with skill contacts
// @Tags Contact
// @Description This function is used to list all contacts with a given skill for an user.
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param skill_name query string true "skill_name"
// @Success 200 {array} db.Contact
// @Router /contacts-with-skill [get]
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

	// Check for owernership.
	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Payload)
	for _, contact := range contacts {
		if contact.Owner != authPayload.Username {
			err := errors.New("contact doesn't belong to the user")
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
	}
	ctx.JSON(http.StatusOK, contacts)
}

// Request holder for get contact with a specific skill and level request.
type getContactWithSkillAndLevelRequest struct {
	SkillName  string `form:"skill_name" binding:"required"`
	SkillLevel string `form:"skill_level" binding:"required"`
}

// getContactWithSkillAndLevel godoc
// @Security bearerAuth
// @Summary Get all the contacts with skill and level contacts
// @Tags Contact
// @Description This function is used to list all contacts with a given skill and a given level for an user.
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param skill_name query string true "skill_name"
// @Param skill_level query string true "skill_level"
// @Success 200 {array} db.Contact
// @Router /contacts-with-skill-and-level [get]
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

	// Check for owernership.
	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Payload)
	for _, contact := range contacts {
		if contact.Owner != authPayload.Username {
			err := errors.New("contact doesn't belong to the user")
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
	}
	ctx.JSON(http.StatusOK, contacts)
}
