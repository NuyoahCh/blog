// @Author Gopher
// @Date 2025/2/11 09:37:00
// @Desc
package controller

import (
	"net/http"

	"blog/service"

	gintemplate "github.com/foolin/gin-template"
	"github.com/gin-gonic/gin"
)

var indexChannelService service.ChannelService
var indexPostService service.PostService

// 后台首页
func AdminIndex(c *gin.Context) {
	gintemplate.HTML(c, http.StatusOK, "index", nil)
}

// 前台首页
func Index(c *gin.Context) {
	channels := indexChannelService.GetChannelList()
	posts := indexPostService.GetPostList()
	gintemplate.HTML(c, http.StatusOK, "index", gin.H{"clist": channels, "posts": posts})
}
