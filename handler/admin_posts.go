package handler

import (
	"net/http"

	"github.com/caris-events/tunalog/config"
	"github.com/caris-events/tunalog/entity"
	"github.com/caris-events/tunalog/store"
	"github.com/gin-gonic/gin"
)

//=======================================================
// View
//=======================================================

func AdminPostsView(c *gin.Context) {
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
	page, err := queryPage(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	countPerPage := 2

	q := &store.ListPostsQuery{
		Offset:        (page - 1) * countPerPage,
		Limit:         countPerPage,
		Title:         c.Query("title"),
		AuthorID:      c.Query("author_id"),
		Visibility:    entity.Visibility(c.Query("visibility")),
		OrderBy:       store.ListPostsOrderBy(c.Query("order_by")),
		PublishedDate: c.Query("published_date"),
		HasPassword:   queryBoolPtr(c, "has_password"),
	}
	posts, err := store.Instance.ListPosts(c, q)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	count, err := store.Instance.CountPosts(c, q)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	dates, err := store.Instance.ListPostDates(c)
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
	c.HTML(http.StatusOK, "admin_posts", gin.H{
		"Config":        config.Instance,
		"Query":         q,
		"IsQuerySetted": q.Title != "" || q.AuthorID != "" || q.Visibility != entity.VisibilityUnknown || q.PublishedDate != "" || q.HasPassword != nil || q.OrderBy != store.ListPostsOrderByCreatedAtDesc,
		"Self":          self,
		"Path":          path(c),
		"Message":       message(c),
		"Posts":         posts,
		"Dates":         dates,
		"Users":         users,
		"Pagination":    p,
	})
}
