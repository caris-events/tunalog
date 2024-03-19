package handler

import (
	"net/http"
	"time"

	"github.com/caris-events/tunalog/entity"
	"github.com/caris-events/tunalog/server/config"
	"github.com/caris-events/tunalog/store"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thanhpk/randstr"
	"golang.org/x/crypto/bcrypt"
)

//=======================================================
// View
//=======================================================

func WizardView(c *gin.Context) {
	if config.Instance != nil {
		c.Redirect(http.StatusFound, "/admin")
		return
	}
	c.HTML(http.StatusOK, "wizard", gin.H{
		"Message": message(c),
	})
}

//=======================================================
// Handle
//=======================================================

type WizardRequest struct {
	Name        string `form:"name"`
	URL         string `form:"url"`
	Description string `form:"description"`
	Username    string `form:"username"`
	Password    string `form:"password"`
	Nickname    string `form:"nickname"`
}

func Wizard(c *gin.Context) {
	if config.Instance != nil {
		c.Redirect(http.StatusFound, "/admin")
		return
	}

	var req *WizardRequest
	if err := c.Bind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	u := &entity.UserW{
		ID:          uuid.New().String(),
		Username:    req.Username,
		Password:    string(hashedPwd),
		Nickname:    req.Nickname,
		Bio:         "",
		AvatarPath:  "",
		CreatedAt:   time.Now().Unix(),
		SuspendedAt: 0,
	}
	if err := store.Instance.CreateUser(c, u); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	conf := &config.Config{
		Name:              req.Name,
		URL:               req.URL,
		Description:       req.Description,
		IsPublic:          true,
		DateFormat:        "2006-01-02",
		TimeFormat:        "15:04:05",
		Timezone:          0,
		InjectedHead:      "",
		InjectedFoot:      "",
		InjectedPostStart: "",
		InjectedPostEnd:   "",
		IsPoweredByShown:  true,
		IsCopyrightShown:  true,
		ColorScheme:       config.ColorSchemeAuto,
		ContainerWidth:    "710px",
		FontFamily:        config.FontFamilyNotoSans,
		FontSize:          "16px",
		HighlightJS:       true,
		AuthorBlock:       config.AuthorBlockStart,
		PostsPerPage:      10,
		AuthSecret:        randstr.String(64),
		Theme:             "default",
	}
	if err := conf.Save(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	setUserID(c, u.ID)
	c.Redirect(http.StatusFound, "/admin/posts")
}
