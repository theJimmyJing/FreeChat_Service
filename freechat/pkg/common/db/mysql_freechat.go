package db

import (
	"Open_IM/pkg/common/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func (m *mysqlDB) FreechatGormDB() (*gorm.DB, error) {
	return m.GormDB(config.FreechatConfig.Mysql.DBAddress[0], config.FreechatConfig.Mysql.DBDatabaseName)
}
