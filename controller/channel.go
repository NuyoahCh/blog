// @Author Gopher
// @Date 2025/2/11 08:48:00
// @Desc 创建操作博客频道controller
package controller

import (
	"blog/models"
	"blog/service"
	"fmt"
	gintemplate "github.com/foolin/gin-template"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var channel service.ChannelService

// ListChannel channel list
func ListChannel(c *gin.Context) {
	c2 := channel.GetChannelList()
	gintemplate.HTML(c, http.StatusOK, "channel/list", gin.H{"clist": c2})
}

// ViewChannel view channel
func ViewChannel(c *gin.Context) {
	sid, r := c.GetQuery("id")
	var chann models.Channel
	if r {
		id, _ := strconv.Atoi(sid)
		chann = channel.GetChannel(id)
	}
	gintemplate.HTML(c, http.StatusOK, "channel/view", gin.H{"channel": chann})
}

// DeleteChannel delete channel
func DeleteChannel(c *gin.Context) {
	sid, _ := c.GetQuery("id")
	id, _ := strconv.Atoi(sid)
	channel.DelChannel(id)
	c.Redirect(http.StatusFound, "/admin/channel/list")
}

// SaveChannel add or update
func SaveChannel(c *gin.Context) {
	var chann models.Channel
	if err := c.ShouldBind(&chann); err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	chann.Status, _ = strconv.Atoi(c.PostForm("status"))

	id, _ := c.GetPostForm("id")

	if id != "0" {
		channel.UpdateChannel(chann)
	} else {
		channel.AddChannel(chann)
	}
	c.Redirect(http.StatusFound, "/admin/channel/list")
}
