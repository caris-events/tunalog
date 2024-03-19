package handler

import (
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/caris-events/tunalog/config"
	"github.com/caris-events/tunalog/translation"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var Instance *Handler

type Handler struct {
	Engine *gin.Engine
}

func init() {
	r := gin.Default()
	Instance = &Handler{r}

	r.Use(sessions.Sessions("tunalog", cookie.NewStore([]byte("hello"))), wizardRedirect)

	funcs := template.FuncMap{
		"str":  translation.Str,
		"strf": translation.Strf,
		"add": func(x, y int) int {
			return x + y
		},
		"sub": func(x, y int) int {
			return x - y
		},
		"unix2date": func(v int64) string {
			return time.Unix(v, 0).Format(config.Instance.DateFormat)
		},
		"timezone": func(v int) string {
			return time.Unix(time.Now().Unix()+int64(v), 0).UTC().Format("2006-01-02 03:04 PM")
		},
	}

	render := multitemplate.NewRenderer()
	render.AddFromFilesFuncs("wizard", funcs, "view/wizard.html")

	r.HTMLRender = render
	r.Static("/assets", "assets")
	r.Static("/files", "files")
	r.Static("/uploads", "uploads")

	r.POST("/wizard", Wizard)
	r.GET("/wizard", WizardView)

}

const (
	KeyMessage      = "message"
	KeyMessageTitle = "message_title"
)

// message reads the flash message from the session, and then deletes it.
func message(c *gin.Context) string {
	session := sessions.Default(c)

	msg, ok := session.Get(KeyMessage).(string)
	if !ok {
		return ""
	}
	session.Delete(KeyMessage)
	session.Save()
	return msg
}

// setMessage writes the flash message to the session.
func setMessage(c *gin.Context, value string) {
	session := sessions.Default(c)

	session.Set(KeyMessage, value)
	session.Save()
}

// setUserID
func setUserID(c *gin.Context, id string) {
	session := sessions.Default(c)

	session.Set("user", id)
	session.Save()
}

// userID
func userID(c *gin.Context) string {
	session := sessions.Default(c)

	id, ok := session.Get("user").(string)
	if !ok {
		return ""
	}
	return id
}

// wizardRedirect
func wizardRedirect(c *gin.Context) {
	if config.Instance != nil {
		c.Next()
		return
	}
	if c.FullPath() != "/wizard" && !strings.Contains(c.FullPath(), "/assets") {
		c.Redirect(http.StatusFound, "/wizard")
		return
	}
	c.Next()
}
