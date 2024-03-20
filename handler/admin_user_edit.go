package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/caris-events/tunalog/config"
	"github.com/caris-events/tunalog/entity"
	"github.com/caris-events/tunalog/store"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//=======================================================
// View
//=======================================================

func AdminUserEditView(c *gin.Context) {
	self, err := self(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	u, err := store.Instance.GetUser(c, c.Param("id"))
	if err != nil {
		setMessage(c, "User not found")
		c.Redirect(http.StatusNotFound, "/admin/users")
		return
	}
	users, err := store.Instance.ListUsers(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	count, err := store.Instance.CountPostsByUser(c, c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "admin_user_edit", gin.H{
		"Config":     config.Instance,
		"Users":      users,
		"User":       u,
		"Self":       self,
		"IsOnlyUser": len(users) == 1,
		"PostCount":  count,
		"Path":       path(c),
		"Message":    message(c),
	})
}

//=======================================================
// Handle
//=======================================================

func AdminUserEdit(c *gin.Context) {
	switch c.PostForm("action") {
	case "reset_avatar":
		AdminUserEditResetAvatar(c)
	case "update":
		AdminUserEditUpdate(c)
	case "delete":
		AdminUserEditDelete(c)
	case "suspend":
		AdminUserEditSuspend(c)
	}
}

//=======================================================
// Reset Avatar
//=======================================================

func AdminUserEditResetAvatar(c *gin.Context) {
	u, err := store.Instance.GetUser(c, c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	w := &entity.UserW{
		ID:          c.Param("id"),
		Username:    u.Username,
		Password:    u.Password,
		Nickname:    u.Nickname,
		Bio:         u.Bio,
		AvatarPath:  "",
		CreatedAt:   u.CreatedAt,
		SuspendedAt: u.SuspendedAt,
	}
	if err := store.Instance.UpdateUser(c, w); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := os.Remove(filepath.FromSlash(u.AvatarPath)); err != nil && !os.IsNotExist(err) {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	setMessage(c, "Avatar reset")
	c.Redirect(http.StatusFound, "/admin/user/"+c.Param("id"))
}

//=======================================================
// Update
//=======================================================

type AdminUserEditRequest struct {
	Username string `form:"username"`
	Password string `form:"password"`
	Nickname string `form:"nickname"`
	Bio      string `form:"bio"`
}

func AdminUserEditUpdate(c *gin.Context) {
	var req *AdminUserEditRequest
	if err := c.Bind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("bind: %w", err))
		return
	}
	targetUID := c.Param("id")

	u, err := store.Instance.GetUser(c, targetUID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	w := &entity.UserW{
		ID:          targetUID,
		Username:    req.Username,
		Password:    u.Password,
		Nickname:    req.Nickname,
		Bio:         req.Bio,
		AvatarPath:  u.AvatarPath,
		CreatedAt:   u.CreatedAt,
		SuspendedAt: u.SuspendedAt,
	}
	if req.Password != "" {
		hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("bcrypt: %w", err))
			return
		}
		w.Password = string(hashedPwd)
	}
	// update the avatar path if a new file is uploaded
	avatarPath, err := saveUpload(c, "avatar_file", photoTypeAvatar)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("save upload: %w", err))
		return
	}
	if avatarPath != "" {
		w.AvatarPath = avatarPath
	}
	if err := store.Instance.UpdateUser(c, w); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	setMessage(c, "user updated")
	c.Redirect(http.StatusFound, fmt.Sprintf("/admin/user/%s", targetUID))
}

//=======================================================
// Delete
//=======================================================

type AdminUserEditDeleteRequest struct {
	TransferToID string `form:"transfer_to_id"`
}

func AdminUserEditDelete(c *gin.Context) {
	var req *AdminUserEditDeleteRequest
	if err := c.Bind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	targetUID := c.Param("id")

	self, err := store.Instance.GetUser(c, userID(c))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if self.ID == targetUID {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	userCount, err := store.Instance.CountUsers(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if userCount == 1 {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if req.TransferToID != "" {
		if _, err := store.Instance.GetUser(c, req.TransferToID); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if err := store.Instance.TransferPosts(c, targetUID, req.TransferToID); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	} else {
		if err := store.Instance.DeletePostsByUser(c, targetUID); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	if err := store.Instance.DeleteUser(c, targetUID); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	setMessage(c, "User deleted")
	c.Redirect(http.StatusFound, "/admin/users")
}

//=======================================================
// Suspend
//=======================================================

type AdminUserEditSuspendRequest struct {
	IsSuspended bool `form:"is_suspended"`
}

func AdminUserEditSuspend(c *gin.Context) {
	var req *AdminUserEditSuspendRequest
	if err := c.Bind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	self, err := store.Instance.GetUser(c, userID(c))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	targetUID := c.Param("id")

	if self.ID == targetUID {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if req.IsSuspended {
		if err := store.Instance.SuspendUser(c, targetUID); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	} else {
		if err := store.Instance.UnsuspendUser(c, targetUID); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	setMessage(c, "User SuspendUser")
	c.Redirect(http.StatusFound, fmt.Sprintf("/admin/user/%s", targetUID))
}
