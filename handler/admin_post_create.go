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

func AdminPostCreateView(c *gin.Context) {
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
	tags, err := store.Instance.ListTags(c, 0, 999, "")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "admin_post_create", gin.H{
		"Config":   config.Instance,
		"Self":     self,
		"Path":     path(c),
		"Message":  message(c),
		"Users":    users,
		"Tags":     tags,
		"RandomID": uuid.New().String(),
	})
}

//=======================================================
// Handle
//=======================================================

type AdminPostCreateRequest struct {
	Title       string            `form:"title"`
	Slug        string            `form:"slug"`
	Excerpt     string            `form:"excerpt"`
	AuthorID    string            `form:"author_id"`
	Password    string            `form:"password"`
	Visibility  entity.Visibility `form:"visibility"`
	Content     string            `form:"content"`
	PublishedAt int64             `form:"published_at"`
	IsPinned    bool              `form:"is_pinned"`
	IsHidden    bool              `form:"is_hidden"`
	Tags        string            `form:"tags"`
}

func AdminPostCreate(c *gin.Context) {
	var req *AdminPostCreateRequest
	if err := c.Bind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	coverPath, err := saveUpload(c, "cover_file", photoTypePostCover)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ids, err := createTags(req.Tags)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	p := &entity.PostW{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Slug:        req.Slug,
		Excerpt:     req.Excerpt,
		CoverPath:   coverPath,
		AuthorID:    req.AuthorID,
		Password:    "",
		Visibility:  req.Visibility,
		Content:     req.Content,
		PublishedAt: req.PublishedAt,
		TagIDs:      ids,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		IsDraft:     c.PostForm("action") == "draft",
		IsHidden:    req.IsHidden,
	}
	if req.PublishedAt == 0 {
		p.PublishedAt = time.Now().Unix()
	}
	if req.IsPinned {
		p.PinnedAt = time.Now().Unix()
	}
	if req.Visibility == entity.VisibilityPassword {
		p.Password = req.Password
	}
	if err := store.Instance.CreatePost(c, p); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	setMessage(c, "Post created")
	c.Redirect(http.StatusSeeOther, "/admin/posts")
}
