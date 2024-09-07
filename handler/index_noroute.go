package handler

import (
	"bytes"
	"net/http"

	"github.com/caris-events/tunalog/store"
	"github.com/caris-events/tunalog/system"
	"github.com/gin-gonic/gin"
)

// ===============================
// NoRouteView
// ===============================

func NoRouteView(c *gin.Context) {
	noRoute(c)
}

func noRoute(c *gin.Context) {
	if system.NotFoundTmpl == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	navs, err := store.ListNavigations()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var tpl bytes.Buffer
	if err := system.NotFoundTmpl.Execute(&tpl, data(c, gin.H{
		"Navigations": navs,
	})); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Data(http.StatusNotFound, "text/html; charset=utf-8", tpl.Bytes())
}
