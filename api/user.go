package api

import (
	"net/http"
	"time"

	"github.com/IsuruHaupe/web-api/auth"
	db "github.com/IsuruHaupe/web-api/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username            string    `json:"username"`
	Fullname            string    `json:"fullname"`
	Email               string    `json:"email"`
	PasswordLastChanged time.Time `json:"password_last_changed"`
	CreateAt            time.Time `json:"create_at"`
}

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

	response := userResponse{
		Username:            user.Username,
		Fullname:            user.Fullname,
		Email:               user.Email,
		PasswordLastChanged: user.PasswordLastChanged,
		CreateAt:            user.CreateAt,
	}
	ctx.JSON(http.StatusOK, response)
}
