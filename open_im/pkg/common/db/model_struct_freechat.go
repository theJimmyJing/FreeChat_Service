package db

type UserFaces struct {
	ID      int    `gorm:"column:id;type:varchar(255)" json:"id"`
	FaceURL string `gorm:"column:face_url;type:varchar(255)" json:"faceURL"`
}
