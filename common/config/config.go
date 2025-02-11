// @Author Gopher
// @Date 2025/2/11 00:14:00
// @Desc
package config

// 组合全部配置模型
type Config struct {
	Server Server `mapstructure:"server"`
	Mysql  Mysql  `mapstructure:"mysql"`
}

// 服务启动端口号配置
type Server struct {
	Post string `mapstructure:"post"`
}

// MySQL数据源配置
type Mysql struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Url      string `mapstructure:"url"`
}
