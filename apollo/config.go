package apollo

type PG struct {
	Path         string `mapstructure:"path"`           // 服务器地址
	Port         string `mapstructure:"port"`           // 服务器端口
	DbName       string `mapstructure:"db-name"`        // 数据库名
	Username     string `mapstructure:"username"`       // 数据库用户名
	Password     string `mapstructure:"password"`       // 数据库密码
	MaxIdleConns int    `mapstructure:"max-idle-conns"` // 空闲中的最大连接数
	MaxOpenConns int    `mapstructure:"max-open-conns"` // 打开到数据库的最大连接数
	LogMode      bool   `mapstructure:"log-mode"`       // 是否开启Gorm全局日志
	LogZap       bool   `mapstructure:"log-zap"`        // 是否通过zap写入日志文件
	LogLevel     string `mapstructure:"log-level"`      // 日志级别
}

type Redis struct {
	DB       int    `mapstructure:"db"`       // redis的哪个数据库
	Addr     string `mapstructure:"addr"`     // 服务器地址:端口
	Password string `mapstructure:"password"` // 密码
}
