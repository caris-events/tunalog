package handler

import (
	"net/http"
	"time"

	"github.com/caris-events/tunalog/entity"
	"github.com/caris-events/tunalog/store"
	"github.com/caris-events/tunalog/system"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/teacat/i18n"
	csrf "github.com/utrack/gin-csrf"
	"golang.org/x/crypto/bcrypt"
)

// ============================
//  WizardView
// ============================

func WizardView(c *gin.Context) {
	if system.Config != nil {
		c.Redirect(http.StatusFound, "admin")
		return
	}
	// Users can manually set the language through param to coordinate with the front-end refresh:w
	locale := c.Query("locale")
	if locale == "" {
		locale = i18n.ParseAcceptLanguage(c.GetHeader("Accept-Language"))[0] // use user's language by default in wizard
	}
	system.ReloadLocale(locale)

	c.HTML(http.StatusOK, "wizard", gin.H{
		"Message":       message(c),
		"Locales":       entity.Locales,
		"DefaultLocale": locale,
		"CSRF":          csrf.GetToken(c),
	})
}

// ============================
//  Wizard
// ============================

type WizardRequest struct {
	Name        string `form:"name" binding:"required,min=1,max=48" conform:"trim"`
	Description string `form:"description" binding:"omitempty,max=128" conform:"trim"`
	Email       string `form:"email" binding:"required,email" conform:"trim"`
	Password    string `form:"password" binding:"required,min=1,max=128"`
	Nickname    string `form:"nickname" binding:"required,min=1,max=32" conform:"trim"`
	Timezone    int    `form:"timezone" binding:"min=-43200,max=50400"`
	Locale      string `form:"locale" binding:"required"`
}

func Wizard(c *gin.Context, req *WizardRequest) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	u := &entity.UserW{
		ID:        uuid.New().String(),
		Email:     req.Email,
		Password:  string(hashedPwd),
		Nickname:  req.Nickname,
		Bio:       "",
		CreatedAt: time.Now().Unix(),
	}
	if err := store.CreateUser(u); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	system.Config = &entity.Config{
		Name:              req.Name,
		Description:       req.Description,
		IsPublic:          true,
		DateFormat:        "2006-01-02",
		TimeFormat:        "15:04",
		Timezone:          req.Timezone,
		InjectedHead:      "",
		InjectedFoot:      "",
		InjectedPostStart: "",
		InjectedPostEnd:   "",
		FooterText:        `Powered by <a href="https://tunalog.org" target="_blank">Tunalog</a> 🐟`,
		ColorScheme:       entity.ColorSchemeAuto,
		ContainerWidth:    "medium",
		FontFamily:        entity.FontFamilyNotoSans,
		FontSize:          "medium",
		HighlightJS:       true,
		AuthorBlock:       entity.AuthorBlockStart,
		PostsPerPage:      10,
		Theme:             "default",
		Locale:            req.Locale,
	}
	if err := system.SaveConfig(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// default post <3
	p := &entity.PostW{
		ID:          uuid.New().String(),
		Title:       system.Locale.String("defaultpost_title"),
		Slug:        "hello-world",
		Excerpt:     "",
		AuthorID:    u.ID,
		Password:    "",
		Visibility:  entity.VisibilityPublic,
		Content:     system.Locale.String("defaultpost_content"),
		PublishedAt: time.Now().Unix() - 60, // so it won't collide with user post created in 1 mins after site initialized
		TagIDs:      []string{},
		CreatedAt:   time.Now().Unix() - 60,
		UpdatedAt:   time.Now().Unix() - 60,
	}
	if err := store.CreatePost(p); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	setUserID(c, u.ID)
	c.Redirect(http.StatusFound, "admin")
}
