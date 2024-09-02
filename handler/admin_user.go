package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/caris-events/tunalog/entity"
	"github.com/caris-events/tunalog/store"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// ============================
//  UsersView
// ============================

func UsersView(c *gin.Context) {
	users, err := store.ListUsers()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "admin_users", data(c, gin.H{
		"Users": users,
	}))
}

// ============================
//  UserCreate
// ============================

type UserCreateRequest struct {
	Email    string `form:"email" binding:"required,email" conform:"trim"`
	Password string `form:"password" binding:"required,min=1,max=128" conform:"trim"`
}

func UserCreate(c *gin.Context, req *UserCreateRequest) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := store.CreateUser(&entity.UserW{
		ID:        uuid.New().String(),
		Email:     req.Email,
		Password:  string(hashedPwd),
		Nickname:  strings.Split(req.Email, "@")[0],
		Bio:       "",
		CreatedAt: time.Now().Unix(),
	}); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	setMessage(c, "notice_user_created")
	c.Redirect(http.StatusFound, "users")
}

// ============================
//  UserEditView
// ============================

func UserEditView(c *gin.Context) {
	u, err := store.GetUser(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	users, err := store.ListUsers()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	count, err := store.CountPostsByUser(u.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "admin_user_edit", data(c, gin.H{
		"Users":      users,
		"User":       u,
		"IsOnlyUser": len(users) == 1,
		"PostCount":  count,
	}))
}

// ============================
//  UserEdit
// ============================

type UserEditRequest struct {
	Email    string `form:"email" binding:"required,email" conform:"trim"`
	Password string `form:"password" binding:"omitempty,min=1,max=128" conform:"trim"`
	Nickname string `form:"nickname" binding:"required,min=1,max=32" conform:"trim"`
	Bio      string `form:"bio" binding:"max=255" conform:"trim"`
}

func UserEdit(c *gin.Context, req *UserEditRequest) {
	_, err := store.GetUser(c.Param("id")) // check if user exists
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if req.Password != "" {
		hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if err := store.UpdateUserPassword(c.Param("id"), string(hashedPwd)); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	if err := store.UpdateUser(c.Param("id"), req.Nickname, req.Bio, req.Email); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	setMessage(c, "notice_user_updated")
	c.Redirect(http.StatusFound, "../users")
}

// ============================
//  UserDelete
// ============================

type UserDeleteRequest struct {
	TransferToID string `form:"transfer_to_id" binding:"omitempty,uuid"`
}

func UserDelete(c *gin.Context, req *UserDeleteRequest) {
	id := c.Param("id")

	if req.TransferToID != "" {
		if _, err := store.GetUser(req.TransferToID); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if err := store.TransferPosts(id, req.TransferToID); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		setMessage(c, "notice_user_deletedwithposts")
	} else {
		if err := store.DeletePostsByUser(id); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		setMessage(c, "notice_user_deleted")
	}
	if err := store.DeleteUser(id); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusFound, "../../users")
}
