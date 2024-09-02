package handler

import (
	"net/http"

	"github.com/caris-events/tunalog/entity"
	"github.com/caris-events/tunalog/system"
	"github.com/gin-gonic/gin"
)

// ===============================
// AppearancesView
// ===============================

func AppearancesView(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_appearances", data(c, gin.H{
		"Themes": system.Themes(),
	}))
}

// ===============================
// AppearancesEdit
// ===============================

type AppearancesEditRequest struct {
	FooterText     string             `form:"footer_text" conform:"trim"`
	ColorScheme    entity.ColorScheme `form:"color_scheme" binding:"omitempty,oneof=light dark"`
	ContainerWidth string             `form:"container_width" binding:"oneof=small medium large"`
	FontFamily     entity.FontFamily  `form:"font_family" binding:"oneof=sans serif"`
	FontSize       string             `form:"font_size" binding:"oneof=small medium large"`
	HighlightJS    bool               `form:"highlight_js"`
	AuthorBlock    entity.AuthorBlock `form:"author_block" binding:"oneof=none start end"`
	PostsPerPage   int                `form:"posts_per_page" binding:"min=1,max=999"`
	Theme          string             `form:"theme" binding:"required"`
}

func AppearancesEdit(c *gin.Context, req *AppearancesEditRequest) {
	if !system.ThemeExists(req.Theme) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	system.Config.FooterText = req.FooterText
	system.Config.ColorScheme = req.ColorScheme
	system.Config.ContainerWidth = req.ContainerWidth
	system.Config.FontFamily = req.FontFamily
	system.Config.FontSize = req.FontSize
	system.Config.HighlightJS = req.HighlightJS
	system.Config.AuthorBlock = req.AuthorBlock
	system.Config.PostsPerPage = req.PostsPerPage
	system.Config.Theme = req.Theme

	if err := system.SaveConfig(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	setMessage(c, "notice_appearances_updated")
	c.Redirect(http.StatusFound, "appearances")
}

// ===============================
// AppearancesEditInjected
// ===============================

type AppearancesEditInjectedRequest struct {
	InjectedHead      string `form:"injected_head" conform:"trim"`
	InjectedFoot      string `form:"injected_foot" conform:"trim"`
	InjectedPostStart string `form:"injected_post_start" conform:"trim"`
	InjectedPostEnd   string `form:"injected_post_end" conform:"trim"`
}

func AppearancesEditInjected(c *gin.Context, req *AppearancesEditInjectedRequest) {
	system.Config.InjectedHead = req.InjectedHead
	system.Config.InjectedFoot = req.InjectedFoot
	system.Config.InjectedPostStart = req.InjectedPostStart
	system.Config.InjectedPostEnd = req.InjectedPostEnd

	if err := system.SaveConfig(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	setMessage(c, "notice_injected_updated")
	c.Redirect(http.StatusFound, "../appearances")
}
