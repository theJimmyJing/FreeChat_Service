package db

const Small = 1
const Large = 2

type UserFaces struct {
	FaceType int8 `gorm:"column:face_type;type:tinyint" json:"face_type"`
	//Platform string `gorm:"column:platform;type:varchar(10)" json:"platform"`
	FaceURL string `gorm:"column:face_url;type:varchar(255)" json:"faceURL"`
}
