package http

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/fs"
	"log"
	"net/http"
	"stzbHelper/http/route/api"
	"stzbHelper/web"
	"strings"
)

func RegisterRoute(r *gin.Engine) {
	staticRoute(r)
	api.Register(r.Group("/v1"))
}

func staticRoute(r *gin.Engine) {
	distFS, err := fs.Sub(web.PublicAssets, "dist")
	if err != nil {
		fmt.Println("初始化静态资源出错")
		return
	}
	teamFS, err := fs.Sub(web.PublicAssets, "team")
	if err != nil {
		fmt.Println("初始化 team 静态资源出错")
		return
	}

	distServer := http.FileServer(http.FS(distFS))
	teamServer := http.FileServer(http.FS(teamFS))

	r.NoRoute(func(c *gin.Context) {
		// 获取请求路径
		reqpath := c.Request.URL.Path

		// 处理根路径默认指向index.html
		if reqpath == "/" {
			reqpath = "/index.html"

		} else if strings.HasSuffix(reqpath, "/") {
			reqpath = strings.TrimSuffix(reqpath, "/")
		}

		tryServe := func(fsys fs.FS, server http.Handler) (served bool) {
			f, err := fsys.Open(strings.TrimPrefix(reqpath, "/"))
			if errors.Is(err, fs.ErrNotExist) {
				return false
			}
			if err != nil {
				log.Println("静态资源访问错误: " + err.Error())
				c.AbortWithStatus(http.StatusInternalServerError)
				return true
			}
			defer f.Close()

			fi, err := f.Stat()
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return true
			}
			if fi.IsDir() {
				c.AbortWithStatus(http.StatusForbidden)
				return true
			}
			server.ServeHTTP(c.Writer, c.Request)
			return true
		}

		if tryServe(distFS, distServer) {
			return
		}

		if reqpath == "/data.html" || reqpath == "/favicon.ico" || strings.HasPrefix(reqpath, "/assets/") {
			if tryServe(teamFS, teamServer) {
				return
			}
		}

		c.JSON(404, gin.H{
			"message": "404 - Page Not Found",
		})
	})
}
