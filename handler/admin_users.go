package handler

import (
	"net/http"
	"time"

	"github.com/caris-events/tunalog/config"
	"github.com/caris-events/tunalog/entity"
	"github.com/caris-events/tunalog/store"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

//=======================================================
// View
//=======================================================

func AdminUsersView(c *gin.Context) {
	self, err := self(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	users, err := store.Instance.ListUsers(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "admin_users", gin.H{
		"Message": message(c),
		"Path":    path(c),
		"Config":  config.Instance,
		"Users":   users,
		"Self":    self,
	})
}

//=======================================================
// Handle
//=======================================================

type AdminUsersRequest struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func AdminUsers(c *gin.Context) {
	var req *AdminUsersRequest
	if err := c.Bind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := store.Instance.CreateUser(c, &entity.UserW{
		ID:          uuid.New().String(),
		Username:    req.Username,
		Password:    string(hashedPwd),
		Nickname:    req.Username,
		Bio:         "",
		AvatarPath:  "",
		CreatedAt:   time.Now().Unix(),
		SuspendedAt: 0,
	}); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	setMessage(c, "User created")

	c.Redirect(http.StatusFound, "/admin/users")
}
