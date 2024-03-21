package handler

import (
	"fmt"
	"net/http"

	"github.com/caris-events/tunalog/config"
	"github.com/caris-events/tunalog/entity"
	"github.com/caris-events/tunalog/store"
	"github.com/gin-gonic/gin"
)

//=======================================================
// View
//=======================================================

func AdminTagEditView(c *gin.Context) {
	self, err := self(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	tag, err := store.Instance.GetTag(c, c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "admin_tag_edit", gin.H{
		"Config":  config.Instance,
		"Self":    self,
		"Path":    path(c),
		"Message": message(c),
		"Tag":     tag,
	})
}

//=======================================================
// Handle
//=======================================================

func AdminTagEdit(c *gin.Context) {
	switch c.PostForm("action") {
	case "update":
		AdminTagEditUpdate(c)
	case "delete":
		AdminTagEditDelete(c)
	}

}

//=======================================================
// Update
//=======================================================

type AdminTagEditUpdateRequest struct {
	Slug        string `form:"slug"`
	Name        string `form:"name"`
	Description string `form:"description"`
}

func AdminTagEditUpdate(c *gin.Context) {
	var req *AdminTagEditUpdateRequest
	if err := c.Bind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	id := c.Param("id")
	tag, err := store.Instance.GetTag(c, id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := store.Instance.UpdateTag(c, &entity.TagW{
		ID:          id,
		Slug:        req.Slug,
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   tag.CreatedAt,
	}); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	setMessage(c, "tag updated")
	c.Redirect(http.StatusFound, fmt.Sprintf("/admin/tag/%s", id))
}

//=======================================================
// Delete
//=======================================================

func AdminTagEditDelete(c *gin.Context) {
	if err := store.Instance.DeleteTag(c.Param("id")); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	setMessage(c, "tag deleted")
	c.Redirect(http.StatusFound, "/admin/tags")
}
