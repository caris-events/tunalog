package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/caris-events/tunalog/config"
	"github.com/caris-events/tunalog/entity"
	"github.com/caris-events/tunalog/store"
	"github.com/gin-gonic/gin"
)

//=======================================================
// View
//=======================================================

func AdminPostEditView(c *gin.Context) {
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
	post, err := store.Instance.GetPost(c, c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	tags, err := store.Instance.ListTags(c, 0, 999, "")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "admin_post_edit", gin.H{
		"Config":  config.Instance,
		"Self":    self,
		"Path":    path(c),
		"Message": message(c),
		"Users":   users,
		"Tags":    tags,
		"Post":    post,
	})
}

//=======================================================
// Handle
//=======================================================

func AdminPostEdit(c *gin.Context) {
	switch c.PostForm("action") {
	case "update", "publish", "draft":

		AdminPostEditUpdate(c)
	case "delete":
		AdminPostEditDelete(c)
	}
}

//=======================================================
// Update
//=======================================================

type AdminPostEditUpdateRequest struct {
	Title        string            `form:"title"`
	Slug         string            `form:"slug"`
	Excerpt      string            `form:"excerpt"`
	AuthorID     string            `form:"author_id"`
	Password     string            `form:"password"`
	Visibility   entity.Visibility `form:"visibility"`
	Content      string            `form:"content"`
	PublishedAt  int64             `form:"published_at"`
	IsPinned     bool              `form:"is_pinned"`
	IsHidden     bool              `form:"is_hidden"`
	IsClearCover bool              `form:"is_clear_cover"`
	Tags         string            `form:"tags"`
}

func AdminPostEditUpdate(c *gin.Context) {
	var req *AdminPostEditUpdateRequest
	if err := c.Bind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	id := c.Param("id")
	post, err := store.Instance.GetPost(c, id)
	if err != nil {
		setMessage(c, "Post not found")
		c.Redirect(http.StatusFound, "/admin/posts")
		return
	}
	ids, err := createTags(req.Tags)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("create tags: %v", err))
		return
	}
	coverPath, err := saveUpload(c, "cover_file", photoTypePostCover)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	p := &entity.PostW{
		ID:          id,
		Title:       req.Title,
		Slug:        req.Slug,
		Excerpt:     req.Excerpt,
		CoverPath:   post.CoverPath,
		AuthorID:    req.AuthorID,
		Visibility:  req.Visibility,
		Content:     req.Content,
		PublishedAt: req.PublishedAt,
		TagIDs:      ids,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   time.Now().Unix(),
		IsHidden:    req.IsHidden,
	}
	if coverPath != "" {
		p.CoverPath = coverPath
	}
	if req.IsClearCover {
		p.CoverPath = ""
	}
	if req.IsPinned {
		p.PinnedAt = time.Now().Unix()
	}
	if req.Visibility == entity.VisibilityPassword {
		p.Password = req.Password
	}
	if c.PostForm("action") == "draft" || (c.PostForm("action") == "update" && post.IsDraft) {
		p.IsDraft = true
	}
	if err := store.Instance.UpdatePost(c, p); err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("update post: %v", err))
		return
	}

	switch c.PostForm("action") {
	case "update":
		setMessage(c, "Post updated")
	case "draft":
		setMessage(c, "Post saved as draft")
	case "publish":
		if time.Now().Unix() < p.PublishedAt {
			setMessage(c, "Post scheduled")
		} else {
			setMessage(c, "Post published")
		}
	}
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/post/%s", id))
}

//=======================================================
// Delete
//=======================================================

func AdminPostEditDelete(c *gin.Context) {
	if err := store.Instance.DeletePost(c, c.Param("id")); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	setMessage(c, "Post deleted")
	c.Redirect(http.StatusSeeOther, "/admin/posts")
}
