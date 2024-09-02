package handler

import (
	"net/http"
	"runtime"
	"time"

	"github.com/caris-events/tunalog/entity"
	"github.com/caris-events/tunalog/system"
	"github.com/gin-gonic/gin"
)

// ============================
//  SettingsView
// ============================

func SettingsView(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_settings", data(c, gin.H{
		"Version":            system.Version,
		"RuntimeVersion":     runtime.Version(),
		"Timezones":          entity.Timezones,
		"Locales":            entity.Locales,
		"IsCustomTimeFormat": system.Config.IsCustomTimeFormat(),
		"IsCustomDateFormat": system.Config.IsCustomDateFormat(),
		"Year":               time.Now().Format("2006"),
		"Month":              time.Now().Format("01"),
		"Day":                time.Now().Format("02"),
		"Hour":               time.Now().Format("03"),
		"Hour24":             time.Now().Format("15"),
		"Minute":             time.Now().Format("04"),
		"Clock":              time.Now().Format("PM"),
	}))
}

// ============================
//  SettingsEdit
// ============================

type SettingsEditRequest struct {
	Name             string `form:"name" binding:"required,max=64" conform:"trim"`
	Description      string `form:"description" binding:"required,max=128" conform:"trim"`
	IsPublic         bool   `form:"is_public"`
	Timezone         int    `form:"timezone" binding:"min=-43200,max=50400"`
	DateFormat       string `form:"date_format" binding:"required"`
	DateFormatCustom string `form:"date_format_custom" conform:"trim"`
	TimeFormat       string `form:"time_format" binding:"required"`
	TimeFormatCustom string `form:"time_format_custom" conform:"trim"`
	Locale           string `form:"locale" binding:"required"`
}

func SettingsEdit(c *gin.Context, req *SettingsEditRequest) {
	system.Config.Name = req.Name
	system.Config.Description = req.Description
	system.Config.IsPublic = req.IsPublic
	system.Config.Timezone = req.Timezone
	system.Config.Locale = req.Locale

	if req.DateFormat == "custom" {
		system.Config.DateFormat = req.DateFormatCustom
	} else {
		system.Config.DateFormat = req.DateFormat
	}
	if req.TimeFormat == "custom" {
		system.Config.TimeFormat = req.TimeFormatCustom
	} else {
		system.Config.TimeFormat = req.TimeFormat
	}
	if err := system.SaveConfig(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	setMessage(c, "notice_settings_updated")
	c.Redirect(http.StatusFound, "settings")
}
