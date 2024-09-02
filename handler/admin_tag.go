package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/caris-events/tunalog/entity"
	"github.com/caris-events/tunalog/store"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ===============================
// TagsView
// ===============================

func TagsView(c *gin.Context) {
	var (
		page         = queryPage(c)
		countPerPage = 50
		keyword      = c.Query("keyword")
	)
	tags, err := store.ListTags((page-1)*countPerPage, countPerPage, keyword)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	count, err := store.CountTags(keyword)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "admin_tags", data(c, gin.H{
		"Keyword":    keyword,
		"Tags":       tags,
		"Pagination": pagination(c, page, count, countPerPage),
	}))
}

// ===============================
// TagCreate
// ===============================

type TagCreateRequest struct {
	Name        string `form:"name" binding:"required,max=64" conform:"trim"`
	Description string `form:"description" binding:"max=255" conform:"trim"`
}

func TagCreate(c *gin.Context, req *TagCreateRequest) {
	if err := store.CreateTag(&entity.TagW{
		ID:          uuid.New().String(),
		Slug:        toSlug(req.Name),
		Name:        strings.ReplaceAll(req.Name, ",", ""),
		Description: req.Description,
		CreatedAt:   time.Now().Unix(),
	}); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	setMessage(c, "notice_tag_created")
	c.Redirect(http.StatusFound, "tags")
}

// ===============================
// TagEditView
// ===============================

func TagEditView(c *gin.Context) {
	tag, err := store.GetTag(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "admin_tag_edit", data(c, gin.H{
		"Tag": tag,
	}))
}

// ===============================
// TagEdit
// ===============================

type TagEditRequest struct {
	Slug        string `form:"slug" binding:"required,max=64" conform:"trim"`
	Name        string `form:"name" binding:"required,max=64" conform:"trim"`
	Description string `form:"description" binding:"max=255" conform:"trim"`
}

func TagEdit(c *gin.Context, req *TagEditRequest) {
	id := c.Param("id")
	tag, err := store.GetTag(id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := store.UpdateTag(&entity.TagW{
		ID:          id,
		Slug:        toSlug(req.Slug),
		Name:        strings.ReplaceAll(req.Name, ",", ""),
		Description: req.Description,
		CreatedAt:   tag.CreatedAt,
	}); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	setMessage(c, "notice_tag_updated")
	c.Redirect(http.StatusFound, "../tags")
}

// ===============================
// TagDelete
// ===============================

func TagDelete(c *gin.Context) {
	if err := store.DeleteTag(c.Param("id")); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	setMessage(c, "notice_tag_deleted")
	c.Redirect(http.StatusFound, "../../tags")
}
