// @Author Gopher
// @Date 2025/2/11 09:03:00
// @Desc 准备路由
package initialize

import (
	"fmt"
	"net/http"

	"blog/common/global"
	"blog/controller"

	gintemplate "github.com/foolin/gin-template"
	"github.com/gin-gonic/gin"
)

func Router() {

	engine := gin.Default()

	// static resource request mapping
	engine.Static("/assets", "./assets")
	engine.StaticFS("/static", http.Dir("./static"))

	engine.HTMLRender = gintemplate.New(gintemplate.TemplateConfig{
		Root:         "templates/frontend",
		Extension:    ".html",
		Master:       "layouts/master",
		DisableCache: true,
	})

	// frontend
	engine.GET("/", controller.Index)
	engine.GET("/post/:id", controller.PostDetail)

	mw := gintemplate.NewMiddleware(gintemplate.TemplateConfig{
		Root:         "templates/backend",
		Extension:    ".html",
		Master:       "layouts/master",
		DisableCache: true,
	})

	// backend admin the interface of frontend
	web := engine.Group("/admin", mw)

	{
		// index
		web.GET("/", controller.AdminIndex)
	}

	{
		// user login API
		web.GET("/channel/list", controller.ListChannel)
		web.GET("/channel/view", controller.ViewChannel)
		web.POST("/channel/save", controller.SaveChannel)
		web.GET("/channel/del", controller.DeleteChannel)
	}

	{
		// post
		web.GET("/post/list", controller.ListPost)
		web.GET("/post/view", controller.ViewPost)
		web.POST("/post/save", controller.SavePost)
		web.GET("/post/del", controller.DeletePost)

		web.POST("/post/upload", controller.UploadThumbnails)

	}

	// run and listen the port
	post := fmt.Sprintf(":%s", global.Config.Server.Post)
	if err := engine.Run(post); err != nil {
		fmt.Printf("server start error: %s", err)
	}
}
