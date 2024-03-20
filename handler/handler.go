package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/caris-events/tunalog/config"
	"github.com/caris-events/tunalog/entity"
	"github.com/caris-events/tunalog/store"
	"github.com/caris-events/tunalog/translation"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sunshineplan/imgconv"
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
	render.AddFromFilesFuncs("admin_users", funcs, "view/admin/_base.html", "view/admin/users.html")
	render.AddFromFilesFuncs("admin_user_edit", funcs, "view/admin/_base.html", "view/admin/user_edit.html")

	r.HTMLRender = render
	r.Static("/assets", "view/assets")
	r.Static("/files", "files")
	r.Static("/uploads", "uploads")

	r.POST("/wizard", Wizard)
	r.GET("/wizard", WizardView)

	r.GET("/admin", AdminView)

	r.GET("/admin/users", AdminUsersView)
	r.POST("/admin/users", AdminUsers)

	r.GET("/admin/user/:id", AdminUserEditView)
	r.POST("/admin/user/:id", AdminUserEdit)

}

const (
	KeyMessage      = "message"
	KeyMessageTitle = "message_title"
	KeyUserID       = "user_id"
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

	session.Set(KeyUserID, id)
	session.Save()
}

// userID
func userID(c *gin.Context) string {
	session := sessions.Default(c)

	id, ok := session.Get(KeyUserID).(string)
	if !ok {
		return ""
	}
	return id
}

func self(c *gin.Context) (*entity.UserR, error) {
	u, err := store.Instance.GetUser(c, userID(c))
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return u, nil
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

func path(c *gin.Context) string {
	switch c.FullPath() {
	case "/admin/users", "/admin/user/:id":
		return "user"
	case "/admin/posts", "/admin/post/create", "/admin/post/:id":
		return "post"
	case "/admin/tags", "/admin/tag/:id":
		return "tag"
	case "/admin/photos":
		return "media"
	case "/admin/navigations":
		return "navigation"
	case "/admin/settings":
		return "settings"
	case "/admin/appearances":
		return "appearances"
	}
	return ""
}

type photoType string

const (
	photoTypePostCover          photoType = "post_cover"
	photoTypePostPhoto          photoType = "post_photo"
	photoTypePostPhotoThumbnail photoType = "post_photo_thumbnail"
	photoTypeAvatar             photoType = "avatar"
)

func saveUpload(c *gin.Context, key string, typ photoType) (string, error) {
	year := time.Now().Format("2006")
	month := time.Now().Format("01")
	id := strconv.Itoa(int(time.Now().Unix())) + "_" + uuid.New().String() // unix timestamp prefix to order by time

	dst := filepath.Join("uploads", year, month, id+".jpg")

	file, err := c.FormFile(key)
	if err != nil {
		return "", nil
	}
	if err := c.SaveUploadedFile(file, dst); err != nil {
		return "", fmt.Errorf("save uploaded file: %w", err)
	}
	srcImg, err := imgconv.Open(dst)
	if err != nil {
		return "", fmt.Errorf("open: %w", err)
	}
	resizeImg := imgconv.Resize(srcImg, &imgconv.ResizeOption{Width: 128})

	dstImg, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}

	if err = imgconv.Write(dstImg, resizeImg, &imgconv.FormatOption{Format: imgconv.JPEG}); err != nil {
		return "", fmt.Errorf("write: %w", err)
	}

	return fmt.Sprintf("uploads/%s/%s/%s.jpg", year, month, id), nil
}
