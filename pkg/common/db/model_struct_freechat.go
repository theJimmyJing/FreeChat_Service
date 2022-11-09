package db

// Male Female 男1女2
const Male = 1
const Female = 2

type UserFaces struct {
	Gender int8 `gorm:"column:gender;type:tinyint" json:"gender"`
	//Platform string `gorm:"column:platform;type:varchar(10)" json:"platform"`
	FaceURL string `gorm:"column:face_url;type:varchar(255)" json:"faceURL"`
}
