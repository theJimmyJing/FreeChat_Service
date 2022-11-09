package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	_ "github.com/jinzhu/gorm"
)

func GetFacesURL() ([]db.UserFaces, error) {
	dbConn, err := db.DB.MysqlDB.FreechatGormDB()
	if err != nil {
		return nil, err
	}
	var r []db.UserFaces
	return r, dbConn.Table("user_faces").Find(&r).Error
}
