package handler

import (
	"fmt"
	"net/http"

	"github.com/caris-events/tunalog/config"
	"github.com/gin-gonic/gin"
)

//=======================================================
// View
//=======================================================

func AdminAppearancesView(c *gin.Context) {
	self, err := self(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "admin_appearances", gin.H{
		"Config":  config.Instance,
		"Self":    self,
		"Path":    path(c),
		"Message": message(c),
	})
}

//=======================================================
// Handle
//=======================================================

func AdminAppearances(c *gin.Context) {
	switch c.PostForm("action") {
	case "update":
		AdminAppearancesUpdate(c)
	case "update_injected":
		AdminAppearancesUpdateInjected(c)
	}
	c.Redirect(http.StatusFound, "/admin/appearances")
}

//=======================================================
// Update
//=======================================================

type AdminAppearancesUpdateRequest struct {
	IsPoweredByShown bool               `form:"is_powered_by_shown"`
	IsCopyrightShown bool               `form:"is_copyright_shown"`
	ColorScheme      config.ColorScheme `form:"color_scheme"`
	ContainerWidth   string             `form:"container_width"`
	FontFamily       config.FontFamily  `form:"font_family"`
	FontSize         string             `form:"font_size"`
	HighlightJS      bool               `form:"highlight_js"`
	AuthorBlock      config.AuthorBlock `form:"author_block"`
	PostsPerPage     int                `form:"posts_per_page"`
}

func AdminAppearancesUpdate(c *gin.Context) {
	var req *AdminAppearancesUpdateRequest
	if err := c.Bind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("bind: %w", err))
		return
	}
	config.Instance.IsPoweredByShown = req.IsPoweredByShown
	config.Instance.IsCopyrightShown = req.IsCopyrightShown
	config.Instance.ColorScheme = req.ColorScheme
	config.Instance.ContainerWidth = req.ContainerWidth
	config.Instance.FontFamily = req.FontFamily
	config.Instance.FontSize = req.FontSize
	config.Instance.HighlightJS = req.HighlightJS
	config.Instance.AuthorBlock = req.AuthorBlock
	config.Instance.PostsPerPage = req.PostsPerPage
	if err := config.Instance.Save(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	setMessage(c, "Saved")
}

//=======================================================
// Update Injected
//=======================================================

type AdminAppearancesUpdateInjectedRequest struct {
	InjectedHead      string `form:"injected_head"`
	InjectedFoot      string `form:"injected_foot"`
	InjectedPostStart string `form:"injected_post_start"`
	InjectedPostEnd   string `form:"injected_post_end"`
}

func AdminAppearancesUpdateInjected(c *gin.Context) {
	var req *AdminAppearancesUpdateInjectedRequest
	if err := c.Bind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("bind: %w", err))
		return
	}
	config.Instance.InjectedHead = req.InjectedHead
	config.Instance.InjectedFoot = req.InjectedFoot
	config.Instance.InjectedPostStart = req.InjectedPostStart
	config.Instance.InjectedPostEnd = req.InjectedPostEnd

	if err := config.Instance.Save(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	setMessage(c, "Saved")
}
