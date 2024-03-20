package handler

import (
	"net/http"
	"time"

	"github.com/caris-events/tunalog/config"
	"github.com/caris-events/tunalog/entity"
	"github.com/caris-events/tunalog/store"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//=======================================================
// View
//=======================================================

func AdminTagsView(c *gin.Context) {
	self, err := self(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	page, err := queryPage(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	countPerPage := 2
	keyword := c.Query("keyword")

	tags, err := store.Instance.ListTags(c, (page-1)*countPerPage, countPerPage, keyword)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	count, err := store.Instance.CountTags(c, keyword)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	p := &entity.Pagination{
		CurrentPage: page,
		TotalCount:  count,
		TotalPages:  totalPages(count, countPerPage),
		Query:       paginationQuery(c),
	}

	c.HTML(http.StatusOK, "admin_tags", gin.H{
		"Config":     config.Instance,
		"Keyword":    c.Query("keyword"),
		"Self":       self,
		"Path":       path(c),
		"Message":    message(c),
		"Tags":       tags,
		"Pagination": p,
	})
}

//=======================================================
// Handle
//=======================================================

type AdminTagsRequest struct {
	Slug        string `form:"slug"`
	Name        string `form:"name"`
	Description string `form:"description"`
}

func AdminTags(c *gin.Context) {
	var req *AdminTagsRequest
	if err := c.Bind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := store.Instance.CreateTag(c, &entity.TagW{
		ID:          uuid.New().String(),
		Slug:        req.Slug,
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now().Unix(),
	}); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	setMessage(c, "tag created")
	c.Redirect(http.StatusFound, "/admin/tags")
}
