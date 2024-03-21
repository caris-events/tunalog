package handler

import (
	"net/http"
	"runtime"
	"time"

	"github.com/caris-events/tunalog/config"
	"github.com/caris-events/tunalog/entity"
	"github.com/gin-gonic/gin"
)

//=======================================================
// View
//=======================================================

func AdminSettingsView(c *gin.Context) {
	self, err := self(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "admin_settings", gin.H{
		"Config":             config.Instance,
		"Self":               self,
		"Path":               path(c),
		"Message":            message(c),
		"Version":            config.Version,
		"RuntimeVersion":     runtime.Version(),
		"Timezones":          entity.Timezones,
		"IsCustomTimeFormat": config.Instance.TimeFormat != "PM 03:04" && config.Instance.TimeFormat != "15:04" && config.Instance.TimeFormat != "03:04 PM",
		"IsCustomDateFormat": config.Instance.DateFormat != "2006-01-02" && config.Instance.DateFormat != "01/02/2006" && config.Instance.DateFormat != "02/01/2006",
		"Year":               time.Now().Format("2006"),
		"Month":              time.Now().Format("01"),
		"Day":                time.Now().Format("02"),
		"Hour":               time.Now().Format("03"),
		"Hour24":             time.Now().Format("15"),
		"Minute":             time.Now().Format("04"),
		"Clock":              time.Now().Format("PM"),
	})
}

//=======================================================
// Handle
//=======================================================

type AdminSettingsRequest struct {
	Name             string `form:"name"`
	Description      string `form:"description"`
	IsPublic         bool   `form:"is_public"`
	URL              string `form:"url"`
	Timezone         int    `form:"timezone"`
	DateFormat       string `form:"date_format"`
	DateFormatCustom string `form:"date_format_custom"`
	TimeFormat       string `form:"time_format"`
	TimeFormatCustom string `form:"time_format_custom"`
}

func AdminSettings(c *gin.Context) {
	var req *AdminSettingsRequest
	if err := c.Bind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	config.Instance.Name = req.Name
	config.Instance.Description = req.Description
	config.Instance.IsPublic = req.IsPublic
	config.Instance.URL = req.URL
	config.Instance.Timezone = req.Timezone

	if req.DateFormat == "custom" {
		config.Instance.DateFormat = req.DateFormatCustom
	} else {
		config.Instance.DateFormat = req.DateFormat
	}
	if req.TimeFormat == "custom" {
		config.Instance.TimeFormat = req.TimeFormatCustom
	} else {
		config.Instance.TimeFormat = req.TimeFormat
	}
	if err := config.Instance.Save(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	setMessage(c, "settings updated")
	c.Redirect(http.StatusFound, "/admin/settings")
}
