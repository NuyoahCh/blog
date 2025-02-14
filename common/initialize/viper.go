// @Author Gopher
// @Date 2025/2/11 00:16:00
// @Desc 使用viper读取配置信息
package initialize

import (
	"blog/common/global"
	"fmt"

	"github.com/spf13/viper"
)

// LoadConfig 加载配置文件
func LoadConfig() {
	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error resources file: %w \n", err))
	}
	if err := viper.Unmarshal(&global.Config); err != nil {
		panic(fmt.Errorf("unable to decode into struct %w \n", err))
	}
}
