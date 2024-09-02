package handler

import (
	"net/http"

	"github.com/caris-events/tunalog/store"
	"github.com/caris-events/tunalog/system"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ===============================
// LoginView
// ===============================

func LoginView(c *gin.Context) {
	self, err := self(c)
	if err != nil && !store.IsNotFound(err) {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if self != nil {
		c.Redirect(http.StatusFound, "../admin/posts")
		return
	}
	c.HTML(http.StatusOK, "login", data(c, gin.H{
		"Name": system.Config.Name,
	}))
}

// ===============================
// Login
// ===============================

type LoginRequest struct {
	Email    string `form:"email" binding:"required,email" conform:"trim"`
	Password string `form:"password" binding:"required,min=1,max=128"`
}

func Login(c *gin.Context, req *LoginRequest) {
	u, err := store.GetUserByEmail(req.Email)
	if err != nil {
		setMessage(c, "notice_login_incorrect")
		c.Redirect(http.StatusFound, "login")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		setMessage(c, "notice_login_incorrect")
		c.Redirect(http.StatusFound, "login")
		return
	}
	setUserID(c, u.ID)
	c.Redirect(http.StatusFound, "../admin/posts")
}

// ===============================
// Logout
// ===============================

func Logout(c *gin.Context) {
	unsetUserID(c)
	setMessage(c, "notice_loggedout")
	c.Redirect(http.StatusFound, "../login")
}
