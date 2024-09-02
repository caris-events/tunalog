package handler

import (
	"net/http"
	"sort"
	"strings"

	"github.com/caris-events/tunalog/entity"
	"github.com/caris-events/tunalog/store"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ===============================
// NavigationsView
// ===============================

func NavigationsView(c *gin.Context) {
	navs, err := store.ListNavigations()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "admin_navigations", data(c, gin.H{
		"Navigations": navs,
	}))
}

// ===============================
// NavigationCreate
// ===============================

type NavigationCreateRequest struct {
	Name string `form:"name" binding:"required,max=64" conform:"trim"`
	URL  string `form:"url" binding:"required,url" conform:"trim"`
}

func NavigationCreate(c *gin.Context, req *NavigationCreateRequest) {
	navs, err := store.ListNavigations()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := store.ClearNavigations(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	for i, n := range navs {
		n.Sequence = i + 1
		if err := store.CreateNavigation(&entity.NavigationW{
			ID:       n.ID,
			Name:     n.Name,
			URL:      n.URL,
			Sequence: n.Sequence,
		}); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	if err := store.CreateNavigation(&entity.NavigationW{
		ID:       uuid.New().String(),
		Name:     req.Name,
		URL:      req.URL,
		Sequence: len(navs) + 1,
	}); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	setMessage(c, "notice_nagivation_created")
	c.Redirect(http.StatusFound, "navigations")
}

// ===============================
// NavigationEdit
// ===============================

type NavigationEditRequest struct {
	Names     []string `form:"name[]" binding:"dive,max=64"`
	URLs      []string `form:"url[]" binding:"dive,url"`
	Sequences []int    `form:"sequence[]" binding:"dive,numeric"`
	IsDeleted []bool   `form:"is_deleted[]"`
}

func NavigationEdit(c *gin.Context, req *NavigationEditRequest) {
	var items []*entity.NavigationW
	for i := range req.Names {
		if req.IsDeleted[i] {
			continue
		}
		items = append(items, &entity.NavigationW{
			ID:       uuid.New().String(),
			Name:     strings.TrimSpace(req.Names[i]),
			URL:      strings.TrimSpace(req.URLs[i]),
			Sequence: req.Sequences[i],
		})
	}
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Sequence < items[j].Sequence
	})
	if err := store.ClearNavigations(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	for i, n := range items {
		n.Sequence = i + 1
		if err := store.CreateNavigation(n); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	setMessage(c, "notice_nagivation_updated")
	c.Redirect(http.StatusFound, "../navigations")
}
