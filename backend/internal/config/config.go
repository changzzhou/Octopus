package config

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	MySQL MySQLConfig
}

type MySQLConfig struct {
	DataSource string
}

func (c *Config) MustNewMySqlConn() sqlx.SqlConn {
	return sqlx.NewMysql(c.MySQL.DataSource)
}
