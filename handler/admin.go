package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//=======================================================
// View
//=======================================================

func AdminView(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "/admin/posts")
}
