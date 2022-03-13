package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/IsuruHaupe/web-api/auth"
	db "github.com/IsuruHaupe/web-api/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Request holder when receiving a create user request.
type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// This is expected returned reponse on succesful creation
type userResponse struct {
	Username            string    `json:"username"`
	Fullname            string    `json:"fullname"`
	Email               string    `json:"email"`
	PasswordLastChanged time.Time `json:"password_last_changed"`
	CreatedAt           time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:            user.Username,
		Fullname:            user.Fullname,
		Email:               user.Email,
		PasswordLastChanged: user.PasswordLastChanged,
		CreatedAt:           user.CreateAt,
	}
}

// createUser godoc
// @Summary Create a new user
// @Description This function is used to create a new user account.
// @Tags user
// @Accept json
// @Produce json
// @Param user body api.createUserRequest true "Create User"
// @Success 200 {object} api.userResponse
// @Router /users [post]
func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	// We verify that the JSON is correct, i.e : all fields are present.
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Compute hashed password.
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// In case of no error, we add the contact in the database.
	args := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		Fullname:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.database.CreateUser(ctx, args)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := newUserResponse(user)
	ctx.JSON(http.StatusOK, response)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	SessionToken          string       `json:"session_token"`
	SessionTokenExpiresAt time.Time    `json:"session_token_expires_at"`
	User                  userResponse `json:"user"`
}

// loginUser godoc
// @Security bearerAuth
// @Summary Login an user
// @Description This function is used to authenticate a user providing the username and password.
// @Tags user
// @Accept json
// @Produce json
// @Param user body api.loginUserRequest true "Login User"
// @Success 200 {object} api.loginUserResponse
// @Router /users/login [post]
func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Check if user is correct.
	user, err := server.database.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = auth.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	sessionToken, sessionPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.SessionDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.database.CreateSession(ctx, db.CreateSessionParams{
		ID:           sessionPayload.ID,
		Username:     user.Username,
		SessionToken: sessionToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    sessionPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := loginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		SessionToken:          sessionToken,
		SessionTokenExpiresAt: sessionPayload.ExpiredAt,
		User:                  newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, response)

}
