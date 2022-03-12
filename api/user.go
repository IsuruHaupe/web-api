package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/IsuruHaupe/web-api/auth"
	db "github.com/IsuruHaupe/web-api/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UserResponse struct {
	Username            string    `json:"username"`
	Fullname            string    `json:"fullname"`
	Email               string    `json:"email"`
	PasswordLastChanged time.Time `json:"password_last_changed"`
	CreatedAt           time.Time `json:"created_at"`
}

func newUserResponse(user db.User) UserResponse {
	return UserResponse{
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
// @Param user body api.CreateUserRequest true "Create User"
// @Success 200 {object} api.UserResponse
// @Router /users [post]
func (server *Server) createUser(ctx *gin.Context) {
	var req CreateUserRequest
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

type LoginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}

// loginUser godoc
// @Security bearerAuth
// @Summary Login an user
// @Description This function is used to authenticate a user providing the username and password.
// @Tags user
// @Accept json
// @Produce json
// @Param user body api.LoginUserRequest true "Login User"
// @Success 200 {object} api.LoginUserResponse
// @Router /users/login [post]
func (server *Server) loginUser(ctx *gin.Context) {
	var req LoginUserRequest
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

	accessToken, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := LoginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, response)

}
