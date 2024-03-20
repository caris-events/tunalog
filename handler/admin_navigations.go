package handler

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/caris-events/tunalog/config"
	"github.com/caris-events/tunalog/entity"
	"github.com/caris-events/tunalog/store"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//=======================================================
// View
//=======================================================

func AdminNavigationsView(c *gin.Context) {
	self, err := self(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	navs, err := store.Instance.ListNavigations(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.HTML(200, "admin_navigations", gin.H{
		"Config":      config.Instance,
		"Self":        self,
		"Path":        path(c),
		"Message":     message(c),
		"Navigations": navs,
	})
}

//=======================================================
// Handle
//=======================================================

func AdminNavigations(c *gin.Context) {
	switch c.PostForm("action") {
	case "update":
		AdminNavigationsUpdate(c)
	case "create":
		AdminNavigationsCreate(c)
	}
	c.Redirect(http.StatusFound, "/admin/navigations")
}

//=======================================================
// Update
//=======================================================

type AdminNavigationsUpdateRequest struct {
	Names     []string `form:"name[]"`
	URLs      []string `form:"url[]"`
	Sequences []int    `form:"sequence[]"`
	IsDeleted []bool   `form:"is_deleted[]"`
}

func AdminNavigationsUpdate(c *gin.Context) {
	var req *AdminNavigationsUpdateRequest
	if err := c.Bind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("bind: %w", err))
		return
	}
	var items []*entity.NavigationW
	for i := range req.Names {
		if req.IsDeleted[i] {
			continue
		}
		items = append(items, &entity.NavigationW{
			ID:       uuid.New().String(),
			Name:     req.Names[i],
			URL:      req.URLs[i],
			Sequence: req.Sequences[i],
		})
	}
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Sequence < items[j].Sequence
	})
	if err := store.Instance.ClearNavigations(c); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	for i, n := range items {
		n.Sequence = i + 1
		if err := store.Instance.CreateNavigation(c, n); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	setMessage(c, "Navigations updated")
}

//=======================================================
// Create
//=======================================================

type AdminNavigationsCreateRequest struct {
	Name string `form:"name"`
	URL  string `form:"url"`
}

func AdminNavigationsCreate(c *gin.Context) {
	var req *AdminNavigationsCreateRequest
	if err := c.Bind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("bind: %w", err))
		return
	}
	navs, err := store.Instance.ListNavigations(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := store.Instance.ClearNavigations(c); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	for i, n := range navs {
		n.Sequence = i + 1
		if err := store.Instance.CreateNavigation(c, &entity.NavigationW{
			ID:       n.ID,
			Name:     n.Name,
			URL:      n.URL,
			Sequence: n.Sequence,
		}); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	if err := store.Instance.CreateNavigation(c, &entity.NavigationW{
		ID:       uuid.New().String(),
		Name:     req.Name,
		URL:      req.URL,
		Sequence: len(navs) + 1,
	}); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	setMessage(c, "Navigation created")
}
