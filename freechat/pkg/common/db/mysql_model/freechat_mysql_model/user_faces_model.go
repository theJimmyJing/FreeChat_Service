package im_mysql_model

import (
	"freechat/pkg/common/db"
	_ "github.com/jinzhu/gorm"
)

func GetFacesURL() ([]db.UserFaces, error) {
	dbConn := db.DB.MysqlDB.DefaultGormDB()
	var r []db.UserFaces
	return r, dbConn.Table("user_faces").Select("id,face_url").Find(&r).Error
}
