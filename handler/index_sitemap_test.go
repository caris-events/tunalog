package handler

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/caris-events/tunalog/entity"
	"github.com/gin-gonic/gin"
)

func TestSitemapURL(t *testing.T) {
	cases := []struct {
		name string
		c    *gin.Context
		post *entity.PostR
		want string
	}{
		{
			name: "success",
			c: &gin.Context{
				Request: &http.Request{
					TLS:  nil,
					Host: "example.com",
				},
			},
			post: &entity.PostR{
				ID:        "1",
				Slug:      "test",
				UpdatedAt: 1600000000,
			},
			want: "http://example.com/post/test",
		},
		{
			name: "success with TLS",
			c: &gin.Context{
				Request: &http.Request{
					TLS:  &tls.ConnectionState{},
					Host: "example.com",
				},
			},
			post: &entity.PostR{
				ID:        "1",
				Slug:      "test",
				UpdatedAt: 1600000000,
			},
			want: "https://example.com/post/test",
		},
		{
			name: "success with localhost",
			c: &gin.Context{
				Request: &http.Request{
					TLS:  nil,
					Host: "localhost",
				},
			},
			post: &entity.PostR{
				ID:        "1",
				Slug:      "test2",
				UpdatedAt: 1600000000,
			},
			want: "http://localhost/post/test2",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := sitemapURL(tt.c, tt.post)
			if got != tt.want {
				t.Errorf("case %v failed: got: %s, want: %s", tt.name, got, tt.want)
			}
		})
	}

}
