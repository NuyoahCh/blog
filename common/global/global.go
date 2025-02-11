// @Author Gopher
// @Date 2025/2/11 00:21:00
// @Desc
package global

import (
	"blog/common/config"

	"gorm.io/gorm"
)

var (
	Config config.Config
	Db     *gorm.DB
)
