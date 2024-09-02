package handler

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

// ===============================
// PhotosView
// ===============================

type PhotoFile struct {
	Year     string
	Month    string
	Filename string
}
type PhotoGroup struct {
	Year      string
	Month     string
	Filenames []string
}

func PhotosView(c *gin.Context) {
	var (
		page         = queryPage(c)
		countPerPage = 100
	)
	var files []*PhotoFile
	if err := filepath.WalkDir("data/uploads/images", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		segments := strings.Split(path, string(os.PathSeparator))
		if len(segments) < 5 {
			return nil
		}
		files = append(files, &PhotoFile{
			Year:     segments[3],
			Month:    segments[4],
			Filename: d.Name(),
		})
		return nil
	}); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	// files were sorted old to new by default, so reverse it
	sort.SliceStable(files, func(i, j int) bool {
		return files[i].Filename > files[j].Filename
	})

	offset := (page - 1) * countPerPage
	if offset+countPerPage > len(files) {
		files = files[offset:]
	} else if offset < len(files) {
		files = files[offset : offset+countPerPage]
	}

	var groups []*PhotoGroup
	for _, file := range files {
		if index := slices.IndexFunc(groups, func(group *PhotoGroup) bool {
			return group.Year == file.Year && group.Month == file.Month
		}); index != -1 {
			groups[index].Filenames = append(groups[index].Filenames, file.Filename)
			continue
		}
		groups = append(groups, &PhotoGroup{
			Year:      file.Year,
			Month:     file.Month,
			Filenames: []string{file.Filename},
		})
	}
	c.HTML(http.StatusOK, "admin_photos", data(c, gin.H{
		"Files":      groups,
		"Pagination": pagination(c, page, len(files), countPerPage),
	}))
}

// ===============================
// PhotoCreate
// ===============================

type PhotoCreateRequest struct {
	PhotoFile *multipart.FileHeader `form:"photo_file"`
}

type PhotoCreateResponse struct {
	Path string `json:"path"`
}

func PhotoCreate(c *gin.Context) {
	var req *PhotoCreateRequest
	if err := c.Bind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	file, err := c.FormFile("photo_file")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	path, err := savePhoto(c, file)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, &PhotoCreateResponse{
		Path: path,
	})
}

// ===============================
// PhotoUpload
// ===============================

type PhotoUploadRequest struct {
	PhotoFiles []*multipart.FileHeader `form:"photo_file[]" binding:"required"`
}

func PhotoUpload(c *gin.Context, req *PhotoUploadRequest) {
	for _, file := range req.PhotoFiles {
		if _, err := savePhoto(c, file); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	setMessage(c, "notice_photo_uploaded")
	c.Redirect(http.StatusFound, "photos")
}

// ===============================
// PhotoDelete
// ===============================

type PhotoDeleteRequest struct {
	Path string `form:"path" binding:"required" conform:"trim"`
}

func PhotoDelete(c *gin.Context, req *PhotoDeleteRequest) {
	basePath := "data/uploads/images"
	absBasePath, err := filepath.Abs(basePath)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	filePath := filepath.Join(basePath, filepath.FromSlash(filepath.Clean(req.Path)))
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !strings.HasPrefix(absPath, absBasePath) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	if err := os.Remove(absPath); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	setMessage(c, "notice_photo_deleted")
	c.Redirect(http.StatusFound, fmt.Sprintf("../photos?%s", c.Request.URL.RawQuery))
}
