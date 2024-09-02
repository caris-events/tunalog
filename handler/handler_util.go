package handler

import (
	"fmt"
	"html/template"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/caris-events/tunalog/entity"
	"github.com/caris-events/tunalog/store"
	"github.com/caris-events/tunalog/system"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sunshineplan/imgconv"
	csrf "github.com/utrack/gin-csrf"
	"golang.org/x/time/rate"
)

var (
	regexpSlug       = regexp.MustCompile(`[^A-Za-z0-9\-._~!$&'()*+,;=\p{L}\p{N}]`)
	throttleLimiters = make(map[string]*rate.Limiter)
)

const (
	KeyMessage      = "message"
	KeyMessageTitle = "message_title"
	KeyUserID       = "user_id"
)

func throttle(c *gin.Context) {
	limiter, ok := throttleLimiters[c.ClientIP()]
	if !ok {
		limiter = rate.NewLimiter(rate.Limit(2), 2)
		throttleLimiters[c.ClientIP()] = limiter
	}
	if limiter.Allow() {
		c.Next()
		return
	}
	c.AbortWithStatus(http.StatusTooManyRequests)
}

func toSlug(v string) string {
	return regexpSlug.ReplaceAllString(strings.ReplaceAll(strings.ToLower(v), " ", "-"), "")
}

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

func setMessage(c *gin.Context, value string) {
	session := sessions.Default(c)

	session.Set(KeyMessage, system.Locale.String(value))
	session.Save()
}

func setUserID(c *gin.Context, id string) {
	session := sessions.Default(c)

	session.Set(KeyUserID, id)
	session.Save()
}

func unsetUserID(c *gin.Context) {
	session := sessions.Default(c)

	session.Delete(KeyUserID)
	session.Save()
}

func userID(c *gin.Context) string {
	session := sessions.Default(c)

	id, ok := session.Get(KeyUserID).(string)
	if !ok {
		return ""
	}
	return id
}

func self(c *gin.Context) (*entity.UserR, error) {
	u, err := store.GetUser(userID(c))
	if err != nil {
		return nil, err
	}
	return u, nil
}

func data(c *gin.Context, data gin.H) gin.H {
	self, err := self(c)
	if err != nil && !store.IsNotFound(err) {
		c.AbortWithError(http.StatusInternalServerError, err)
		return nil
	}
	suffix := "https://"
	if c.Request.TLS == nil {
		suffix = "http://"
	}
	data["Self"] = self
	data["Config"] = system.Config
	data["Message"] = message(c)
	data["CSRF"] = csrf.GetToken(c)
	data["URL"] = map[string]string{
		"Root":         filepath.Clean(suffix + c.Request.Host + c.Request.URL.Path + entity.RelativeRoots[c.FullPath()]),
		"Absolute":     suffix + c.Request.Host + c.Request.URL.Path,
		"RelativeRoot": entity.RelativeRoots[c.FullPath()],
		"PageType":     entity.PageTypes[c.FullPath()],
	}
	return data
}

func checkConfig(c *gin.Context) {
	if system.Config == nil {
		c.Redirect(http.StatusFound, "/wizard")
		c.Abort()
		return
	}
	c.Next()
}

func checkPublic(c *gin.Context) {
	if system.Config != nil && !system.Config.IsPublic && userID(c) == "" {
		setMessage(c, "notice_site_private")
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}
	c.Next()
}

func checkLoggedIn(c *gin.Context) {
	if userID(c) == "" {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}
	c.Next()
}

// queryPage gets the page from the query string,
// or returns 1 if not found.
func queryPage(c *gin.Context) int {
	if i, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil {
		return i
	}
	return 1
}

func totalPages(totalItems, itemsPerPage int) int {
	return int(math.Ceil(float64(totalItems) / float64(itemsPerPage)))
}

func paginationQuery(c *gin.Context) template.URL {
	query := c.Request.URL.Query()
	delete(query, "page")

	var newQueryParams []string
	for key, values := range query {
		newQueryParams = append(newQueryParams, fmt.Sprintf("%s=%s", key, values[0]))
	}
	newQueryString := strings.Join(newQueryParams, "&")

	if newQueryString != "" {
		newQueryString += "&"
	}
	return template.URL(newQueryString)
}

func saveCover(c *gin.Context, pid string) (string, error) {
	var (
		localDst  = fmt.Sprintf("data/uploads/covers/%s.jpg", pid)
		publicDst = fmt.Sprintf("uploads/covers/%s.jpg", pid)
	)
	file, err := c.FormFile("cover_file")
	if err != nil {
		return "", nil
	}
	if err := c.SaveUploadedFile(file, localDst); err != nil {
		return "", err
	}
	srcImg, err := imgconv.Open(localDst)
	if err != nil {
		return "", err
	}
	resizeImg := imgconv.Resize(srcImg, &imgconv.ResizeOption{Width: 1024})

	dstImg, err := os.OpenFile(localDst, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return "", err
	}
	if err = imgconv.Write(dstImg, resizeImg, &imgconv.FormatOption{Format: imgconv.JPEG}); err != nil {
		return "", err
	}
	return publicDst, nil
}

func savePhoto(c *gin.Context, file *multipart.FileHeader) (string, error) {
	var (
		year      = time.Now().Format("2006")
		month     = time.Now().Format("01")
		unix      = strconv.Itoa(int(time.Now().Unix()))
		id        = uuid.New().String()
		localDst  = fmt.Sprintf("data/uploads/images/%s/%s/%s_%s.jpg", year, month, unix, id)
		publicDst = fmt.Sprintf("uploads/images/%s/%s/%s_%s.jpg", year, month, unix, id)
	)
	if err := c.SaveUploadedFile(file, localDst); err != nil {
		return "", err
	}
	srcImg, err := imgconv.Open(localDst)
	if err != nil {
		return "", err
	}
	resizeImg := imgconv.Resize(srcImg, &imgconv.ResizeOption{Width: 2000})

	dstImg, err := os.OpenFile(localDst, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return "", err
	}
	if err = imgconv.Write(dstImg, resizeImg, &imgconv.FormatOption{Format: imgconv.JPEG}); err != nil {
		return "", err
	}
	return publicDst, nil
}

func createTags(tagNames string) (ids []string, err error) {
	var names []string
	for _, v := range strings.Split(tagNames, ",") {
		if v = strings.TrimSpace(v); v != "" {
			names = append(names, v)
		}
	}
	tags, err := store.GetTagsByName(names)
	if err != nil {
		return nil, err
	}
	for _, name := range names {
		if index := slices.IndexFunc(tags, func(tag *entity.TagR) bool {
			return name == tag.Name
		}); index != -1 {
			ids = append(ids, tags[index].ID)
			continue
		}
		t := &entity.TagW{
			ID:        uuid.New().String(),
			Slug:      toSlug(name),
			Name:      name,
			CreatedAt: time.Now().Unix(),
		}
		if err := store.CreateTag(t); err != nil {
			return nil, err
		}
		ids = append(ids, t.ID)
	}
	return ids, nil
}

func pagination(c *gin.Context, page, total, countPerPage int) *entity.Pagination {
	return &entity.Pagination{
		CurrentPage: page,
		TotalCount:  total,
		TotalPages:  totalPages(total, countPerPage),
		Query:       paginationQuery(c),
	}
}
