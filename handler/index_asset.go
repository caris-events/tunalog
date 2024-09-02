package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/caris-events/tunalog/system"
	"github.com/gin-gonic/gin"
)

// ===============================
// AssetView
// ===============================

func AssetView(c *gin.Context) {
	theme := "default"
	if system.Config != nil {
		theme = system.Config.Theme
	}
	basePath := fmt.Sprintf("data/themes/%s/assets", theme)

	absBasePath, err := filepath.Abs(basePath)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	asset := c.Param("asset")
	filePath := filepath.Join(basePath, filepath.FromSlash(filepath.Clean(asset)))

	resolvedPath, err := filepath.EvalSymlinks(filePath)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	absResolvedPath, err := filepath.Abs(resolvedPath)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !strings.HasPrefix(absResolvedPath, absBasePath) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	c.File(absResolvedPath)
}
