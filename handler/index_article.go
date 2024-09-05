package handler

import (
	"bytes"
	"net/http"
	"time"

	"github.com/caris-events/tunalog/entity"
	"github.com/caris-events/tunalog/store"
	"github.com/caris-events/tunalog/system"
	"github.com/gin-gonic/gin"
)

// ===============================
// IndexView
// ===============================

type IndexQuery struct {
	Tag    string
	Author string
	Date   string
}

func (q *IndexQuery) IsEmpty() bool {
	return q.Tag == "" && q.Author == "" && q.Date == ""
}

func IndexView(c *gin.Context) {
	self, err := self(c)
	if err != nil && !store.IsNotFound(err) {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var (
		page         = queryPage(c)
		countPerPage = system.Config.PostsPerPage
		query        = &IndexQuery{}
	)
	q := &store.ListPostsQuery{
		Offset:      (page - 1) * countPerPage,
		Limit:       countPerPage,
		Title:       c.Query("title"),
		IsPublished: store.PtrBool(true),
		IsTrashed:   store.PtrBool(false),
	}
	if self == nil {
		q.Visibilities = []entity.Visibility{entity.VisibilityPublic, entity.VisibilityPassword}
	} else {
		q.Visibilities = []entity.Visibility{entity.VisibilityPublic, entity.VisibilityPassword, entity.VisibilityPrivate}
	}
	// tag
	if v := c.Param("tag"); v != "" {
		tag, err := store.GetTagBySlug(v)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		q.TagID = tag.ID
		query.Tag = tag.Name
	}
	// author
	if v := c.Param("author"); v != "" {
		user, err := store.GetUser(v)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		q.AuthorID = user.ID
		query.Author = user.Nickname
	}
	// dates
	if y := c.Param("year"); y != "" {
		q.PublishedYear = y
		query.Date = y

		if m := c.Param("month"); m != "" {
			q.PublishedMonth = m
			query.Date += "/" + m

			if d := c.Param("day"); d != "" {
				q.PublishedDay = d
				query.Date += "/" + d
			}
		}
	}
	posts, err := store.ListPosts(q)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	navs, err := store.ListNavigations()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	count, err := store.CountPosts(q)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var tpl bytes.Buffer
	if err := system.IndexTmpl.Execute(&tpl, data(c, gin.H{
		"Posts":       posts,
		"Pagination":  pagination(c, page, count, countPerPage),
		"Navigations": navs,
		"Filter":      query,
	})); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", tpl.Bytes())
}

// ===============================
// SingularView
// ===============================

func SingularView(c *gin.Context) {
	self, err := self(c)
	if err != nil && !store.IsNotFound(err) {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	p, err := store.GetPostBySlug(c.Param("slug"))
	if err != nil && !store.IsNotFound(err) {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if p == nil {
		c.HTML(http.StatusNotFound,"404NotFound",gin.H{
			"title": "Page Not Found",
		})
		return
	}
	if self == nil && p.Visibility != entity.VisibilityPublic && p.Visibility != entity.VisibilityPassword {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if self == nil && p.PublishedAt > time.Now().Unix() {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	navs, err := store.ListNavigations()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	prevPost, err := store.GetPreviousPost(p.ID)
	if err != nil && !store.IsNotFound(err) {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	nextPost, err := store.GetNextPost(p.ID)
	if err != nil && !store.IsNotFound(err) {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var isUnlocked bool
	if self != nil || p.Visibility == entity.VisibilityPublic {
		isUnlocked = true
	} else {
		if c.PostForm("password") == p.Password {
			isUnlocked = true
		} else if c.Request.Method == http.MethodPost {
			setMessage(c, "notice_post_incorrect")
		}
	}
	var tpl bytes.Buffer
	if err := system.SingularTmpl.Execute(&tpl, data(c, gin.H{
		"Post":         p,
		"Navigations":  navs,
		"PreviousPost": prevPost,
		"NextPost":     nextPost,
		"IsUnlocked":   isUnlocked,
	})); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", tpl.Bytes())
}
